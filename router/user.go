package router

import (
	"database/sql"
	"example/web-service-gin/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserInfoByID(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		userId, err1 := strconv.Atoi(c.Param("userid"))

		if err1 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  "invalid user id",
			})
			return
		}

		info, err2 := database.GetUserInfoByID(db, userId)

		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, info)
	}
}

func ChangeUsername(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var User struct {
			NewUsername string `json:"username" binding:"required"`
		}

		err := c.ShouldBindJSON(&User)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Request body is malformed",
			})
			return
		}

		if len(User.NewUsername) <= 0 || len(User.NewUsername) > 15 {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  "Must be at most 15 characters",
			})
			return
		}

		userId, err1 := strconv.Atoi(c.Param("userid"))

		if err1 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  "invalid user id",
			})
			return
		}

		err2 := database.ChangeUsername(db, userId, User.NewUsername)

		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}
