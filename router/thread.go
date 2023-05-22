package router

import (
	"database/sql"
	"example/web-service-gin/database"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Still a bit messy, sql.DB should not be exposed
// outside of database pkg. However, sufficient for now.
func GetThreadById(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		threadid := c.Param("threadid")
		thread := database.GetThreadById(db, threadid)
		c.JSON(http.StatusOK, gin.H{
			"thread": thread,
		})
	}
}

func EditThreadById(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		threadid, err := strconv.Atoi(c.Param("threadid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Thread id is malformed",
			})
			return
		}

		var EditedThread struct {
			Title   *string `json:"title"`
			Content *string `json:"content"`
			Tags    []int   `json:"tags"`
			UserId  int     `json:"userId"`
			Token   string  `json:"token"`
		}

		err = c.ShouldBindJSON(&EditedThread)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Request body is malformed",
			})
			return
		}

		err2 := database.EditThreadById(db, EditedThread.Title, EditedThread.Content,
			EditedThread.Tags, threadid, EditedThread.UserId, EditedThread.Token)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
		}

		//err := database.EditThreadById(db, threadid)
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}

func PostComment(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		threadid, err := strconv.Atoi(c.Param("threadid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Thread id is malformed",
			})
			return
		}

		var Comment struct {
			AuthorId int    `json:"authorid" binding:"required"`
			Content  string `json:"content" binding:"required"`
			ParentId int    `json:"parentid" binding:"required"`
		}
		err = c.ShouldBindJSON(&Comment)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Request body is malformed",
			})
			return
		}
		fmt.Println("TEST")
		commentId, err2 := database.PostComment(db, Comment.AuthorId, Comment.Content, Comment.ParentId, threadid)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
		}

		//err := database.EditThreadById(db, threadid)
		c.JSON(http.StatusOK, gin.H{
			"status":    "success",
			"commentid": commentId,
		})
	}
}

func DeleteThreadById(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		threadId, err := strconv.Atoi(c.Param("threadid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Thread id is malformed",
			})
			return
		}

		var userData struct {
			Token  string `json:"token" binding:"required"`
			UserId int    `json:"userid" binding:"required"`
		}

		err = c.ShouldBindJSON(&userData)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "JSON body is malformed",
			})
			return
		}

		err = database.DeleteThread(db, threadId, userData.Token, userData.UserId)

		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}
