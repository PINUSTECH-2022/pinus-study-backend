package main

import (
	"example/web-service-gin/database"
	"example/web-service-gin/middlewares"
	"example/web-service-gin/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

// https://seefnasrul.medium.com/create-your-first-go-rest-api-with-jwt-authentication-in-gin-framework-dbe5bda72817

func main() {
	r := gin.Default()

	db := database.GetDb()

	r.Use(middlewares.CORSMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/login", router.LogIn(db))

	r.POST("/signup", router.SignUp(db))

	r.POST("/verify_email/:userid", router.MakeVerification(db))
	r.PUT("/verify_email/:emailid", router.VerifyEmail(db))

	r.POST("/forgot_password/:userid", router.ForgotPassword(db))

	r.POST("/password_recovery/:recoveryid", router.CheckPasswordRecovery(db))

	r.POST("/me", middlewares.JwtAuthMiddleware(), router.GetPersonalInfo(db))

	r.GET("/user/:userid", router.GetUserInfoByID(db))

	r.POST("/user/change_password/:userid", middlewares.JwtAuthMiddleware(), router.ChangePassword(db))

	r.POST("/module", router.GetModules(db))
	r.GET("/module/:moduleid", router.GetModuleByModuleId(db))
	r.POST("/module/:moduleid", middlewares.JwtAuthMiddleware(), router.PostThread(db))

	r.GET("/comment/:id", router.GetCommentById(db))
	r.DELETE("/comment/:id", middlewares.JwtAuthMiddleware(), router.DeleteCommentById(db))
	r.PUT("/comment/:id", middlewares.JwtAuthMiddleware(), router.UpdateCommentById(db))

	r.GET("/thread/:threadid", router.GetThreadById(db))
	r.PUT("/thread/:threadid", router.EditThreadById(db))
	r.POST("/thread/:threadid", middlewares.JwtAuthMiddleware(), router.PostComment(db))
	r.DELETE("/thread/:threadid", router.DeleteThreadById(db))

	r.GET("/bookmark/user/:userid", middlewares.JwtAuthMiddleware(), router.GetBookmark(db))
	r.GET("/bookmark/:threadid/:userid", middlewares.JwtAuthMiddleware(), router.GetBookmarkThread(db))
	r.POST("/bookmark/:threadid/:userid", middlewares.JwtAuthMiddleware(), router.BookmarkThread(db))
	r.DELETE("/bookmark/:threadid/:userid", middlewares.JwtAuthMiddleware(), router.UnbookmarkThread(db))

	r.GET("/subscribes/:moduleid", router.GetSubscribers(db))
	r.GET("/subscribes/:moduleid/:userid", router.DoesSubscribe(db))
	r.POST("/subscribes/:moduleid/:userid", middlewares.JwtAuthMiddleware(), router.Subscribe(db))
	r.DELETE("/subscribes/:moduleid/:userid", middlewares.JwtAuthMiddleware(), router.Unsubscribe(db))

	r.GET("/likes/thread/:threadid/:userid", router.GetLikeThread(db))
	r.POST("/likes/thread/:threadid/:userid/:state", middlewares.JwtAuthMiddleware(), router.SetLikeThread(db))
	r.GET("/likes/comment/:commentid/:userid", router.GetLikeComment(db))
	r.POST("/likes/comment/:commentid/:userid/:state", middlewares.JwtAuthMiddleware(), router.SetLikeComment(db))

	r.GET("/review/:moduleid", router.GetReviewByModule(db))
	r.POST("/review/:moduleid", middlewares.JwtAuthMiddleware(), router.PostReview(db))
	r.GET("/review/:moduleid/:userid", router.GetReviewByModuleAndUser(db))
	r.PUT("/review/:moduleid/:userid", middlewares.JwtAuthMiddleware(), router.EditReviewByModuleAndUser(db))
	r.DELETE("/review/:moduleid/:userid", middlewares.JwtAuthMiddleware(), router.DeleteReviewByModuleAndUser(db))
	r.GET("/review/workload", router.GetWorkload(db))
	r.GET("/review/grade", router.GetGrade(db))
	r.GET("/review/difficulty", router.GetDifficulty(db))
	r.GET("/review/semester", router.GetSemester(db))

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
