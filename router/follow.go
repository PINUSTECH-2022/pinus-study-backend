package router

import (
	"database/sql"
	"example/web-service-gin/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FollowUser(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("userid"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Follow request failed",
			})
			return
		}

		var Follow struct {
			FollowingId int `json:"followingid" binding:"required"`
		}

		err = c.ShouldBindJSON(&Follow)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Request body is malformed",
			})
			return
		}

		err2 := database.FollowUser(db, userId, Follow.FollowingId)
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
