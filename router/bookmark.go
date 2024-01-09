package router

import (
	"database/sql"
	"example/web-service-gin/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Bookmark a thread
func BookmarkThread(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		threadId, err := strconv.Atoi(c.Param("threadid"))
		userId, err1 := strconv.Atoi(c.Param("userid"))
		if err != nil || err1 != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Request body is malformed",
			})
			return
		}

		err2 := database.BookmarkThread(db, threadId, userId)
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

// Unbookmark a thread
func UnbookmarkThread(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		threadId, err := strconv.Atoi(c.Param("threadid"))
		userId, err1 := strconv.Atoi(c.Param("userid"))
		if err != nil || err1 != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Request body is malformed",
			})
			return
		}

		err2 := database.UnbookmarkThread(db, threadId, userId)
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

// Get whether a thread is being bookmarked by the user
func GetBookmarkThread(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		threadId, err := strconv.Atoi(c.Param("threadid"))
		userId, err1 := strconv.Atoi(c.Param("userid"))
		if err != nil || err1 != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Request body is malformed",
			})
			return
		}

		isBookmarked, err2 := database.GetBookmarkThread(db, threadId, userId)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"bookmarked": isBookmarked,
		})
	}
}

func GetBookmark(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("userid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Request body is malformed",
			})
			return
		}

		bookmarkedThreads := database.GetBookmark(db, userId)

		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"threads": bookmarkedThreads,
		})
	}
}