package router

import (
	"database/sql"
	"example/web-service-gin/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeleteCommentBody struct {
	UserId int
	Token  string
}

type UpdateCommentBody struct {
	Content string
	UserId  int
	Token   string
}

func GetCommentById(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause": "Comment id is malformed",
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
				"cause": "Comment id is malformed",
			})
			return
		}

		var requestBody DeleteCommentBody
		bodyErr := c.BindJSON(&requestBody)
		if bodyErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause": "Request body is malformed",
			})
			return
		}

		status := database.DeleteCommentById(db, id, requestBody.UserId,
			requestBody.Token)

		if !status {
			c.JSON(http.StatusNotAcceptable, "fail")
			return
		}

		c.JSON(http.StatusOK, "nice")
	}
}

func UpdateCommentById(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, convErr := strconv.Atoi(c.Param("id"))
		if convErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause": "Comment id is malformed",
			})
			return
		}

		var requestBody UpdateCommentBody
		bodyErr := c.BindJSON(&requestBody)
		if bodyErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause": "Request body is malformed",
			})
		}

		status := database.UpdateCommentById(db, id, requestBody.Content,
			requestBody.UserId, requestBody.Token)

		if !status {
			c.JSON(http.StatusNotAcceptable, "fail")
			return
		}

		c.JSON(http.StatusOK, "nice")

	}
}
