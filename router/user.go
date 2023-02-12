package router

import (
	"database/sql"
	"example/web-service-gin/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserInfoByID(db* sql.DB) func(c * gin.Context) {
	return func (c * gin.Context) {
		userId, err1 := strconv.Atoi(c.Param("userid"))

		if err1 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause": "invalid user id",
			})
			return
		}

		info, err2 := database.GetUserInfoByID(db, userId)
		
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, info)
	}
}
