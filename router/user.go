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

func SignUp(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var User struct {
			Email    string `json:"email"`
			Username string `json:"username"`
			Password string `json:"password"`
		}
		err := c.ShouldBindJSON(&User)
		if err != nil {
			panic(err)
		}

		is_alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(User.Username)
		
		if !is_alphanumeric {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause": "username must be alphanumeric",
			})
			return
		}

		if !isEmailValid(User.Email) {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause": "email is not valid",
			})
			return
		}

		token, err2 := database.SignUp(db, User.Email, User.Username, User.Password)
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
		})
	}
}

func LogIn(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var User struct {
			NameOrEmail string `json:"username"`
			Password string `json:"password"`
		}
		err := c.ShouldBindJSON(&User)
		if err != nil {
			panic(err)
		}

		success, token, err2 := database.LogIn(db, User.NameOrEmail, User.Password)
		
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			panic(err2)
		}

		var status string
		if success {
			status = "success"
		} else {
			status = "failure"
			token = ""
		}
		//err := database.EditThreadById(db, threadid)
		c.JSON(http.StatusOK, gin.H{
			"status": status,
			"token":  token,
		})
	}
}
