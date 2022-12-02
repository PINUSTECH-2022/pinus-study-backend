package main

import (
	"example/web-service-gin/database"
	"example/web-service-gin/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db := database.GetDb()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/module", router.GetModules(db))
	r.GET("/comment/:id", router.GetCommentById(db))

	r.DELETE("/comment/:id", router.DeleteCommentById(db))

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
