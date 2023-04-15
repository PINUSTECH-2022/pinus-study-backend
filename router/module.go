package router

import (
	"database/sql"
	"example/web-service-gin/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Still a bit messy, sql.DB should not be exposed
// outside of database pkg. However, sufficient for now.
func GetModules(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var SearchQuery struct {
			Keyword string `json:"keyword"`
			Page    int    `json:"page"`
		}

		err := c.ShouldBindJSON(&SearchQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Request body is malformed",
			})
			return
		}

		// If page is not specified, default is 1
		if SearchQuery.Page == 0 {
			SearchQuery.Page = 1
		}

		fmt.Println("Search Query Keyword: ", SearchQuery.Keyword)

		modules := database.GetModules(db, SearchQuery.Keyword, SearchQuery.Page)
		c.JSON(http.StatusOK, gin.H{
			"module_list": modules,
		})
	}
}

func GetModuleByModuleId(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		moduleid := c.Param("moduleid")
		fmt.Println("1")
		module := database.GetModuleByModuleId(db, moduleid)
		fmt.Println("2")
		c.JSON(http.StatusOK, gin.H{
			"module": module,
		})
	}
}

func PostThread(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		moduleid := c.Param("moduleid")

		var Module struct {
			AuthorId int    `json:"authorid" binding:"required"`
			Content  string `json:"content" binding:"required"`
			Title    string `json:"title" binding:"required"`
			Tags     []int  `json:"tags" binding:"required"`
		}
		err := c.ShouldBindJSON(&Module)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Request body is malformed",
			})
			return
		}

		threadId, err2 := database.PostThread(db, Module.AuthorId, Module.Content, Module.Title, Module.Tags, moduleid)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
		}

		//err := database.EditThreadById(db, threadid)
		c.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"threadid": threadId,
		})
	}
}
