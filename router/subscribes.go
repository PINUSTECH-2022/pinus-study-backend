package router

import (
	"database/sql"
	"example/web-service-gin/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Still a bit messy, sql.DB should not be exposed
// outside of database pkg. However, sufficient for now.
func GetSubcsribers(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		moduleid := c.Param("moduleid")
		subcsribers := database.GetSubcsribers(db, moduleid)
		c.JSON(http.StatusOK, gin.H{
			"users": subcsribers,
		})
	}
}

func Subcsribe(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		moduleid := c.Param("moduleid")
		userid, err := strconv.Atoi(c.Param("userid"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err.Error(),
			})
			panic(err)
		}

		database.Subcsribe(db, moduleid, userid)

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}

func Unsubcsribe(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		moduleid := c.Param("moduleid")
		userid, err := strconv.Atoi(c.Param("userid"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err.Error(),
			})
			panic(err)
		}

		database.Unsubcsribe(db, moduleid, userid)
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}
