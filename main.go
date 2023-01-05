package main

import (
	"example/web-service-gin/database"
	"example/web-service-gin/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()

	db := database.GetDb()

	r.Use(CORSMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/user", router.LogIn(db))
	r.POST("/user", router.SignUp(db))

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
