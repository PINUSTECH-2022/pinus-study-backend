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
			panic(err)
		}

		var EditedThread struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		err = c.ShouldBindJSON(&EditedThread)
		if err != nil {
			panic(err)
		}

		err2 := database.EditThreadById(db, EditedThread.Title, EditedThread.Content, threadid)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			panic(err2)
		}

		//err := database.EditThreadById(db, threadid)
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}
