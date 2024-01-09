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
		userId, emailId, isEmailExist, isUsernameExist, err1 := database.SignUp(db, User.Email, User.Username, User.Password, secretCode)
		if err1 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err1.Error(),
			})
			return
		}

		if isEmailExist {
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

		var status string
		if success {
			status = "success"
		} else if !isSignedUp {
			status = "username or email does not exist"
		} else if !isVerified {
			status = "failure due to unverified email"
		} else {
			status = "failure due to wrong password"
			userid = -1
		}

		c.JSON(http.StatusOK, gin.H{
			"status": status,
			"token":  token,
			"userid": userid,
		})
	}
}

// Sends the email verification link to the user's email
func sendVerification(userid int, emailid int, email string, username string, secretCode string) error {
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}

	frontendUrl := os.Getenv("FRONTEND_URL")
	emailSenderName := os.Getenv("EMAIL_SENDER_NAME")
	emailSenderAddress := os.Getenv("EMAIL_SENDER_ADDRESS")
	emailSenderPassword := os.Getenv("EMAIL_SENDER_PASSWORD")

	subject := "Welcome to PINUS STUDY!"
	verifyUrl := fmt.Sprintf("%s/verify_email?email_id=%d&secret_code=%s", frontendUrl, emailid, secretCode)
	content := fmt.Sprintf(`Dear Pinusian, <br/>
	There has been a request to register the address %s with the user %s on the PINUS STUDY. 
	In order to complete the address registration you need to go to the following link in a web browser: <a href = "%s">%s</a> <br/>
	Best regards from PINUS`, email, username, verifyUrl, verifyUrl)
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
