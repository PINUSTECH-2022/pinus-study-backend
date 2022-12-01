package router

import (
	"database/sql"
	"example/web-service-gin/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Still a bit messy, sql.DB should not be exposed
// outside of database pkg. However, sufficient for now.
func GetModules(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		modules := database.GetModules(db)
		c.JSON(http.StatusOK, modules)
	}
}
