package router

import (
	"database/sql"
	"example/web-service-gin/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCommentById(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		comment := database.GetCommentById(db, id)

		// if content is empty, there is no comment
		if comment.Content == "" {
			c.JSON(http.StatusNotFound, "")
			return
		}

		c.JSON(http.StatusOK, comment)
	}
}
