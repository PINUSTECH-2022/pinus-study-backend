package router

import (
	"database/sql"
	"example/web-service-gin/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetLikeThread(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		threadid, err := strconv.Atoi(c.Param("threadid"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err.Error(),
			})
			return
		}

		userid, err2 := strconv.Atoi(c.Param("userid"))
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
		}

		state, err3 := database.GetLikeThread(db, threadid, userid)
		if err3 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err3.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"state": state,
		})
	}
}

func SetLikeThread(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		threadid, err := strconv.Atoi(c.Param("threadid"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err.Error(),
			})
			return
		}

		userid, err2 := strconv.Atoi(c.Param("userid"))
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
		}

		state, err3 := strconv.Atoi(c.Param("state"))
		if err3 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err3.Error(),
			})
			return
		}

		if state != 0 && state != 1 && state != -1 {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  "provided state must be in [-1, 0, 1]",
			})
			return
		}

		err4 := database.SetLikeThread(db, state, threadid, userid)
		if err4 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err4.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}

func GetLikeComment(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		commentid, err := strconv.Atoi(c.Param("commentid"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err.Error(),
			})
			return
		}

		userid, err2 := strconv.Atoi(c.Param("userid"))
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
		}

		state, err3 := database.GetLikeComment(db, commentid, userid)
		if err3 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err3.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"state": state,
		})
	}
}

func SetLikeComment(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		commentid, err := strconv.Atoi(c.Param("commentid"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err.Error(),
			})
			return
		}

		userid, err2 := strconv.Atoi(c.Param("userid"))
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
		}

		state, err3 := strconv.Atoi(c.Param("state"))
		if err3 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err3.Error(),
			})
			return
		}

		if state != 0 && state != 1 && state != -1 {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  "provided state must be in [-1, 0, 1]",
			})
			return
		}

		err4 := database.SetLikeComment(db, state, commentid, userid)
		if err4 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err4.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}

func GetLikeReview(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		reviewid, err := strconv.Atoi(c.Param("reviewid"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err.Error(),
			})
			return
		}

		userid, err2 := strconv.Atoi(c.Param("userid"))
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
		}

		state, err3 := database.GetLikeReview(db, reviewid, userid)
		if err3 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err3.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"state": state,
		})
	}
}

func SetLikeReview(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		reviewid, err := strconv.Atoi(c.Param("reviewid"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err.Error(),
			})
			return
		}

		userid, err2 := strconv.Atoi(c.Param("userid"))
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err2.Error(),
			})
			return
		}

		state, err3 := strconv.Atoi(c.Param("state"))
		if err3 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err3.Error(),
			})
			return
		}

		if state != 0 && state != 1 && state != -1 {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  "provided state must be in [-1, 0, 1]",
			})
			return
		}

		err4 := database.SetLikeReview(db, state, reviewid, userid)
		if err4 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err4.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}
