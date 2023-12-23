package router

import (
	"database/sql"
	"example/web-service-gin/database"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func isEmailAvailable(db *sql.DB, email string) bool {
	rows, _ := db.Query("SELECT email FROM Users WHERE email = $1", email)
	for rows.Next() {
		return false
	}
	return true
}

func isUsernameAvailable(db *sql.DB, username string) bool {
	rows, _ := db.Query("SELECT username FROM Users WHERE username = $1", username)
	for rows.Next() {
		return false
	}
	return true
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

		if !isEmailAvailable(db, User.Email) {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":	"email is already taken",
			})
			return
		}

		if !isUsernameAvailable(db, User.Username) {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":	"username is already taken",
			})
			return
		}

		userId, token, err2 := database.SignUp(db, User.Email, User.Username, User.Password)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"token":  token,
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

		success, userid, token, err2 := database.LogIn(db, User.NameOrEmail, User.Password)

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
