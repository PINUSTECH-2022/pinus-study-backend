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
		var SearchQuery struct {
			Keyword string `json:"keyword"`
			Page    int    `json:"page"`
		}

		err := c.ShouldBindJSON(&SearchQuery)
		if err != nil {
			panic(err)
		}

		// If page is not specified, default is 1
		if SearchQuery.Page == 0 {
			SearchQuery.Page = 1
		}

		modules := database.GetModules(db, SearchQuery.Keyword, SearchQuery.Page)
		c.JSON(http.StatusOK, gin.H{
			"module_list": modules,
		})
	}
}

func GetModuleByModuleId(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		moduleid := c.Param("moduleid")
		module := database.GetModuleByModuleId(db, moduleid)
		c.JSON(http.StatusOK, gin.H{
			"module": module,
		})
	}
}

func PostThread(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		moduleid := c.Param("moduleid")

		var Module struct {
			AuthorId int    `json:"authorid"`
			Content  string `json:"content"`
			Title    string `json:"title"`
		}
		err := c.ShouldBindJSON(&Module)
		if err != nil {
			panic(err)
		}

		err2 := database.PostThread(db, Module.AuthorId, Module.Content, moduleid, Module.Title)
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
