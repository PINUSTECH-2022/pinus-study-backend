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

func UnfollowUser(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("userid"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Unfollow request failed",
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

		err2 := database.UnfollowUser(db, userId, Follow.FollowingId)
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

func GetFollowers(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("userid"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Unfollow request failed",
			})
			return
		}

		followers, err := database.GetFollowers(db, userId)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "success",
			"followers": followers,
		})
	}
}

func GetFollowings(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("userid"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Unfollow request failed",
			})
			return
		}

		following, err := database.GetFollowings(db, userId)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"followings": following,
		})
	}
}

// Get a list of thread id which is posted by the user's following
func GetFollowingsThread(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("userid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Followings' thread request failed",
			})
			return
		}

		threads, err1 := database.GetFollowingsThreads(db, userId)
		if err1 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err1.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"threads": threads,
		})
	}
}
