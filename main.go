package main

import (
	"database/sql"
	"example/web-service-gin/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Still a bit messy, sql.DB should not be exposed
// outside of database pkg. However, sufficient for now.
func getModules(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		modules := database.GetModules(db)
		c.JSON(http.StatusOK, modules)
	}
}

func main() {
	r := gin.Default()

	db := database.GetDb()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/module", getModules(db))

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
