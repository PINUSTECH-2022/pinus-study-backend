package router

import (
	"database/sql"
	"example/web-service-gin/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPersonalInfo(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var User struct {
			UserId int `json:"userid"`
		}

		err := c.ShouldBindJSON(&User)
		if err != nil {
			panic(err)
		}

		info, err2 := database.GetPersonalInfo(db, User.UserId)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			panic(err2)
		}
		c.JSON(http.StatusOK, info)

	}

}
