package router

import (
	"database/sql"
	"example/web-service-gin/database"
	"example/web-service-gin/mail"
	"example/web-service-gin/util"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func SignUp(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var User struct {
			Email    string `json:"email" binding:"required"`
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		err := c.ShouldBindJSON(&User)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Request body is malformed",
			})
			return
		}

		is_alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(User.Username)

		if !is_alphanumeric {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  "username must be alphanumeric",
			})
			return
		}

		if !isEmailValid(User.Email) {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  "email is not valid",
			})
			return
		}

		secretCode := util.RandomString(32)
		userId, emailId, isEmailExist, isUsernameExist, isVerified, err1 := database.SignUp(db, User.Email, User.Username, User.Password, secretCode)
		if err1 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err1.Error(),
			})
			return
		}

		if isEmailExist && isVerified {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  "email is already taken",
			})
			return
		}

		if isUsernameExist {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  "username is already taken",
			})
			return
		}

		if isEmailExist && !isVerified {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  "email has been registered but not verified",
				"userid": userId,
			})
			return
		}

		err2 := sendVerification(userId, emailId, User.Email, User.Username, secretCode)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  "verification link can not be sent via email",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"userid": userId,
		})
	}
}

func LogIn(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var User struct {
			NameOrEmail string `json:"username" binding:"required"`
			Password    string `json:"password" binding:"required"`
		}
		err := c.ShouldBindJSON(&User)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Request body is malformed",
			})
			return
		}

		success, isSignedUp, isVerified, userid, token, err2 := database.LogIn(db, User.NameOrEmail, User.Password)

		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
		}

		var status, cause string
		if success {
			status = "success"
		} else {
			status = "failure"
			userid = -1
		}

		if !isSignedUp {
			cause = "username or email does not exist"
		} else if !isVerified {
			cause = "failure due to unverified email"
		} else {
			cause = "failure due to wrong password"
		}

		if status == "success" {
			c.JSON(http.StatusOK, gin.H{
				"status": status,
				"token":  token,
				"userid": userid,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": status,
				"cause":  cause,
			})
		}
	}
}

// Sends the email verification link to the user's email
func sendVerification(userid int, emailid int, email string, username string, secretCode string) error {
	frontendUrl := os.Getenv("FRONTEND_URL")
	emailSenderName := os.Getenv("EMAIL_SENDER_NAME")
	emailSenderAddress := os.Getenv("EMAIL_SENDER_ADDRESS")
	emailSenderPassword := os.Getenv("EMAIL_SENDER_PASSWORD")

	if frontendUrl == "" {
		err := godotenv.Load(".env")

		if err != nil {
			panic(err)
		}

		frontendUrl = os.Getenv("FRONTEND_URL")
		emailSenderName = os.Getenv("EMAIL_SENDER_NAME")
		emailSenderAddress = os.Getenv("EMAIL_SENDER_ADDRESS")
		emailSenderPassword = os.Getenv("EMAIL_SENDER_PASSWORD")
	}

	subject := "Welcome to PINUS STUDY!"
	verifyUrl := fmt.Sprintf("%s/verify_email?email_id=%d&secret_code=%s", frontendUrl, emailid, secretCode)
	content := fmt.Sprintf(`Dear Pinusian, <br/>
	There has been a request to register the address %s with the user %s on the PINUS STUDY. 
	In order to complete the address registration you need to go to the following link in a web browser: <a href = "%s">%s</a> <br/>
	Best regards from PINUS Team`, email, username, verifyUrl, verifyUrl)
	to := []string{email}

	err1 := mail.NewGmailSender(emailSenderName, emailSenderAddress, emailSenderPassword).SendEmail(subject, content, to, nil, nil, nil)
	if err1 != nil {
		return err1
	}

	return nil
}

func VerifyEmail(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		emailId, err := strconv.Atoi(c.Param("emailid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Email id is malformed",
			})
			return
		}

		var User struct {
			SecretCode string `json:"secretcode" binding:"required"`
		}

		err1 := c.ShouldBindJSON(&User)
		if err1 != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "JSON request body is malformed",
			})
			return
		}

		isVerified, isExpired, isMatch, err2 := database.VerifyEmail(db, emailId, User.SecretCode)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
		}

		if isVerified {
			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"message": "Email has already been verified before",
			})
			return
		}

		// Check whether the secret code is expired
		if isExpired {
			c.JSON(http.StatusForbidden, gin.H{
				"status": "failure",
				"cause":  "Verification link is expired",
			})
			return
		}

		// Check whether the secret code in the request match with the secret code in the database
		if !isMatch {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "failure",
				"cause":  "Secret code does not match",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Email has been verified successfully",
		})
	}
}

// Makes the email verification (used when the user ask to make a new verification)
func MakeVerification(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("userid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "userid is malformed",
			})
			return
		}

		secretCode := util.RandomString(32)

		isExist, emailId, email, username, isVerified, err1 := database.StoreSecretCode(db, userId, secretCode)
		if err1 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err1.Error(),
			})
			return
		}

		if !isExist {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "user id not exist",
			})
			return
		}

		if isVerified {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  "email has been verified",
			})
			return
		}

		err2 := sendVerification(userId, emailId, email, username, secretCode)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  "verification link can not be sent via email",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"id":     userId,
		})
	}
}

// Changes the user's password
func ChangePassword(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("userid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "userid is malformed",
			})
			return
		}

		var Password struct {
			New string `json:"newPassword" binding:"required"`
			Old string `json:"oldPassword" binding:"required"`
		}

		err1 := c.ShouldBindJSON(&Password)
		if err1 != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Request body is malformed",
			})
			return
		}

		isUserExist, isVerified, isPasswordMatch, err2 := database.ChangePassword(db, userId, Password.Old, Password.New)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
		}

		if !isUserExist {
			c.JSON(http.StatusNotFound, gin.H{
				"status": "failure",
				"cause":  "User not found",
			})
			return
		}

		if !isVerified {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "failure",
				"cause":  "Email has not been verified",
			})
			return
		}

		if !isPasswordMatch {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "failure",
				"cause":  "Password does not match",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}

// Sends the email verification link to the user's email

func sendPasswordRecovery(recoveryId int, email string, secretCode string) error {
	frontendUrl := os.Getenv("FRONTEND_URL")
	emailSenderName := os.Getenv("EMAIL_SENDER_NAME")
	emailSenderAddress := os.Getenv("EMAIL_SENDER_ADDRESS")
	emailSenderPassword := os.Getenv("EMAIL_SENDER_PASSWORD")

	if frontendUrl == "" {
		err := godotenv.Load(".env")

		if err != nil {
			panic(err)
		}

		frontendUrl = os.Getenv("FRONTEND_URL")
		emailSenderName = os.Getenv("EMAIL_SENDER_NAME")
		emailSenderAddress = os.Getenv("EMAIL_SENDER_ADDRESS")
		emailSenderPassword = os.Getenv("EMAIL_SENDER_PASSWORD")
	}

	subject := "Password Reset Request: Action Required for Your PINUS STUDY Account"
	recoveryUrl := fmt.Sprintf("%s/password_recovery?recovery_id=%d&secret_code=%s", frontendUrl, recoveryId, secretCode)
	content := fmt.Sprintf(`Dear Pinusian, <br/>
	We noticed that there has been a recent request to reset the password for your PINUS STUDY account. No worries, we're here to assist you in regaining access to your account. <br/>
	To initiate the password reset process, please follow the link provided below:<br/>
	<a href = "%s">%s</a> <br/>
	If you didn't request this password reset or if you believe this was an error, please disregard this email. 
	Your account security is important to us, and we recommend taking precautionary measures such as updating your password regularly and enabling two-factor authentication.<br/>
	<br/>
	Best regards,<br/>
	PINUS STUDY Team`, recoveryUrl, recoveryUrl)
	to := []string{email}

	err1 := mail.NewGmailSender(emailSenderName, emailSenderAddress, emailSenderPassword).SendEmail(subject, content, to, nil, nil, nil)
	if err1 != nil {
		return err1
	}

	return nil
}

// Forgot password to create password recovery and send it via email
func ForgotPassword(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var User struct {
			Email string `json:"email" binding:"required"`
		}

		err := c.ShouldBindJSON(&User)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "JSON request body is malformed",
			})
			return
		}

		secretCode := util.RandomString(32)
		isExist, isVerified, recoveryId, email, err1 := database.MakePasswordRecovery(db, User.Email, secretCode)
		if err1 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err1.Error(),
			})
			return
		}

		if !isExist {
			c.JSON(http.StatusNotFound, gin.H{
				"status": "failure",
				"cause":  "User not found",
			})
			return
		}

		if !isVerified {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "failure",
				"cause":  "Email has not been verified",
			})
			return
		}

		err2 := sendPasswordRecovery(recoveryId, email, secretCode)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  "Failure in sending email",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}

// Check whether the password recovery link is valid
func CheckPasswordRecovery(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		recoverId, err := strconv.Atoi(c.Param("recoveryid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Recovery id is malformed",
			})
			return
		}

		var User struct {
			SecretCode string `json:"secretcode" binding:"required"`
		}

		err1 := c.ShouldBindJSON(&User)
		if err1 != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "JSON request body is malformed",
			})
			return
		}

		isMatch, isExpired, isUsed, err2 := database.GetRecoverPassword(db, recoverId, User.SecretCode)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
		}

		// Check whether the secret code in the request match with the secret code in the database
		if !isMatch {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "failure",
				"cause":  "Secret code does not match",
			})
			return
		}

		// Check whether the secret code is expired
		if isExpired {
			c.JSON(http.StatusForbidden, gin.H{
				"status": "failure",
				"cause":  "Recovery link is expired",
			})
			return
		}

		// Check whether the secret code is used
		if isUsed {
			c.JSON(http.StatusGone, gin.H{
				"status": "failure",
				"cause":  "Recovery link has been used",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}

// Password recovery
func RecoverPassword(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		recoverId, err := strconv.Atoi(c.Param("recoveryid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Recovery id is malformed",
			})
			return
		}

		var User struct {
			Password   string `json:"password" binding:"required"`
			SecretCode string `json:"secretcode" binding:"required"`
		}
		err1 := c.ShouldBindJSON(&User)
		if err1 != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "JSON request body is malformed",
			})
			return
		}

		isExist, isExpired, isMatch, isUsed, err2 := database.RecoverPassword(db, recoverId, User.SecretCode, User.Password)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
		}

		// Check whether the recovery code exist
		if !isExist {
			c.JSON(http.StatusNotFound, gin.H{
				"status": "failure",
				"cause":  "Recovery code does not exist",
			})
		}

		// Check whether the secret code in the request match with the secret code in the database
		if !isMatch {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "failure",
				"cause":  "Secret code does not match",
			})
			return
		}

		// Check whether the secret code is expired
		if isExpired {
			c.JSON(http.StatusForbidden, gin.H{
				"status": "failure",
				"cause":  "Recovery link is expired",
			})
			return
		}

		// Check whether the secret code is used
		if isUsed {
			c.JSON(http.StatusGone, gin.H{
				"status": "failure",
				"cause":  "Recovery link has been used",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}
