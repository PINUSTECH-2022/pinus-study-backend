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

// Get the list of user who likes a certain thread
func GetListOfLikeThread(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		threadid, err := strconv.Atoi(c.Param("threadid"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err.Error(),
			})
			return
		}

		listOfLike, err1 := database.GetListOfLikeThread(db, threadid)
		if err1 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err1.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"likes":  listOfLike,
		})
	}
}

// Get the list of user who dislikes a certain thread
func GetListOfDislikeThread(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		threadid, err := strconv.Atoi(c.Param("threadid"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err.Error(),
			})
			return
		}

		listOfDislike, err1 := database.GetListOfDislikeThread(db, threadid)
		if err1 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err1.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"dislikes": listOfDislike,
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

// Get the list of user who likes a certain comment
func GetListOfLikeComment(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		commentid, err := strconv.Atoi(c.Param("commentid"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err.Error(),
			})
			return
		}

		listOfLike, err1 := database.GetListOfLikeComment(db, commentid)
		if err1 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err1.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"likes":  listOfLike,
		})
	}
}

// Get the list of user who dislikes a certain comment
func GetListOfDislikeComment(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		commentid, err := strconv.Atoi(c.Param("commentid"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err.Error(),
			})
			return
		}

		listOfDislike, err1 := database.GetListOfDislikeComment(db, commentid)
		if err1 != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err1.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"dislikes": listOfDislike,
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
