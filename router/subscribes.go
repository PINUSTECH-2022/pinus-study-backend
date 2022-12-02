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
func GetSubscribers(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		moduleid := c.Param("moduleid")
		subcsribers := database.GetSubscribers(db, moduleid)
		c.JSON(http.StatusOK, gin.H{
			"users": subcsribers,
		})
	}
}

func DoesSubscribe(db *sql.DB) func(c *gin.Context) {
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

		res, err2 := database.DoesSubscribe(db, moduleid, userid)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			panic(err2)
		}
		c.JSON(http.StatusOK, gin.H{
			"subscribed": res,
		})
	}
}

func Subscribe(db *sql.DB) func(c *gin.Context) {
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

		err2 := database.Subscribe(db, moduleid, userid)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			panic(err2)
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}

func Unsubscribe(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		moduleid := c.Param("moduleid")
		userid, err := strconv.Atoi(c.Param("userid"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err.Error(),
			})
			return
		}

		err2 := database.Unsubscribe(db, moduleid, userid)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			panic(err2)
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}
