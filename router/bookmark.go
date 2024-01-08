package router

import (
	"database/sql"
	"example/web-service-gin/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func BookmarkThread(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var Bookmark struct {
			ThreadId string `json:"threadid" binding:"required"`
			UserId   string `json:"userid" binding:"required"`
		}

		err := c.ShouldBindJSON(&Bookmark)
		threadId, err1 := strconv.Atoi(Bookmark.ThreadId)
		userId, err2 := strconv.Atoi(Bookmark.UserId)
		if err != nil || err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Request body is malformed",
			})
			return
		}

		err3 := database.BookmarkThread(db, threadId, userId)
		if err3 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err3.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}
