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
	r.GET("/module/:moduleid", router.GetModuleByModuleId(db))
	r.POST("/module/:moduleid", router.PostThread(db))

	r.GET("/comment/:id", router.GetCommentById(db))
	r.DELETE("/comment/:id", router.DeleteCommentById(db))
	r.PUT("/comment/:id", router.UpdateCommentById(db))

	r.GET("/thread/:threadid", router.GetThreadById(db))
	r.PUT("/thread/:threadid", router.EditThreadById(db))
	r.POST("/thread/:threadid", router.PostComment(db))

	r.GET("/subscribes/:moduleid", router.GetSubscribers(db))
	r.GET("/subscribes/:moduleid/:userid", router.DoesSubscribe(db))
	r.POST("/subscribes/:moduleid/:userid", router.Subscribe(db))
	r.DELETE("/subscribes/:moduleid/:userid", router.Unsubscribe(db))

	r.GET("/likes/thread/:threadid/:userid", router.GetLikeThread(db))
	r.POST("/likes/thread/:threadid/:userid/:state", router.SetLikeThread(db))
	r.GET("/likes/comment/:commentid/:userid", router.GetLikeComment(db))
	r.POST("/likes/comment/:commentid/:userid/:state", router.SetLikeComment(db))

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
