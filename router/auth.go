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

		if !database.IsEmailAvailable(db, User.Email) {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  "email is already taken",
			})
			return
		}

		if !database.IsUsernameAvailable(db, User.Username) {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  "username is already taken",
			})
			return
		}

		userId, err2 := database.SignUp(db, User.Email, User.Username, User.Password)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
		}

		err3 := makeVerification(db, userId, User.Email, User.Username)

		if err3 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err3.Error(),
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

		success, isVerified, userid, token, err2 := database.LogIn(db, User.NameOrEmail, User.Password)

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
		} else if !isVerified {
			status = "failure due to unverified email"
		} else {
			status = "failure due to wrong password"
			token = ""
			userid = -1
		}

		c.JSON(http.StatusOK, gin.H{
			"status": status,
			"token":  token,
			"userid": userid,
		})
	}
}

func makeVerification(db *sql.DB, userid int, email string, username string) error {
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}

	frontendUrl := os.Getenv("FRONTEND_URL")
	emailSenderName := os.Getenv("EMAIL_SENDER_NAME")
	emailSenderAddress := os.Getenv("EMAIL_SENDER_ADDRESS")
	emailSenderPassword := os.Getenv("EMAIL_SENDER_PASSWORD")

	secretCode := util.RandomString(32)
	id, err := database.StoreSecretCode(db, userid, email, secretCode)
	if err != nil {
		return err
	}

	subject := "Welcome to PINUS STUDY!"
	verifyUrl := fmt.Sprintf("%s/verify_email?email_id=%d&secret_code=%s", frontendUrl, id, secretCode)
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
