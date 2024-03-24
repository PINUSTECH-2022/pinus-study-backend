package router

import (
	"database/sql"
	"example/web-service-gin/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCommentById(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Comment id is malformed",
			})
			return
		}
		comment := database.GetCommentById(db, id)

		// if content is empty, there is no comment
		if comment.Content == "" {
			c.JSON(http.StatusNotFound, "")
			return
		}

		c.JSON(http.StatusOK, comment)
	}
}

func DeleteCommentById(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, convErr := strconv.Atoi(c.Param("id"))
		if convErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Comment id is malformed",
			})
			return
		}

		var DeleteCommentBody struct {
			UserId int `json:"userid" binding:"required"`
		}

		bodyErr := c.BindJSON(&DeleteCommentBody)
		if bodyErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Request body is malformed",
			})
			return
		}

		err := database.DeleteCommentById(db, id, DeleteCommentBody.UserId)

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

func UpdateCommentById(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, convErr := strconv.Atoi(c.Param("id"))
		if convErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Comment id is malformed",
			})
			return
		}

		var UpdateCommentBody struct {
			UserId  int    `json:"userid" binding:"required"`
			Content string `json:"content" binding:"required"`
		}

		bodyErr := c.BindJSON(&UpdateCommentBody)
		if bodyErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Request body is malformed",
			})
		}

		err := database.UpdateCommentById(db, id, UpdateCommentBody.UserId, UpdateCommentBody.Content)

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

func GetReviewCommentById(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Comment id is malformed",
			})
			return
		}
		comment := database.GetReviewCommentById(db, id)

		// if content is empty, there is no comment
		if comment.Content == "" {
			c.JSON(http.StatusNotFound, "")
			return
		}

		c.JSON(http.StatusOK, comment)
	}
}

func DeleteReviewCommentById(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, convErr := strconv.Atoi(c.Param("id"))
		if convErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Comment id is malformed",
			})
			return
		}

		var DeleteCommentBody struct {
			UserId int `json:"userid" binding:"required"`
		}

		bodyErr := c.BindJSON(&DeleteCommentBody)
		if bodyErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Request body is malformed",
			})
			return
		}

		err := database.DeleteReviewCommentById(db, id, DeleteCommentBody.UserId)

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
