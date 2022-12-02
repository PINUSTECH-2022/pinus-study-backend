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

func GetCommentById(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			panic(err)
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
			panic(convErr)
		}

		var requestBody DeleteCommentBody
		bodyErr := c.BindJSON(&requestBody)
		if bodyErr != nil {
			panic(bodyErr)
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
