package router

import (
	"database/sql"
	"example/web-service-gin/database"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetReviewByModule(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		moduleId := strings.ToUpper(c.Param("moduleid"))

		review := database.GetReviewByModule(db, moduleId)
		c.JSON(http.StatusOK, gin.H{
			"review": review,
		})
	}
}

func PostReview(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		moduleId := strings.ToUpper(c.Param("moduleid"))

		var Review struct {
			UserId        int    `json:"user_id" binding:"required"`
			Workload      string `json:"workload" binding:"required"`
			ExpectedGrade string `json:"expected_grade"`
			ActualGrade   string `json:"actual_grade"`
			Difficulty    string `json:"difficulty" binding:"required"`
			SemesterTaken string `json:"semester_taken" binding:"required"`
			Lecturer      string `json:"lecturer"`
			Content       string `json:"content" binding:"required"`
			Suggestion    string `json:"suggestion"`
		}
		bodyErr := c.ShouldBindJSON(&Review)
		if bodyErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Request body is malformed",
			})
			return
		}

		workload, convErr := strconv.Atoi(Review.Workload)
		if convErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Workload is malformed",
			})
			return
		}

		difficulty, convErr := strconv.Atoi(Review.Difficulty)
		if convErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Difficulty is malformed",
			})
			return
		}

		err := database.PostReview(db, moduleId, Review.UserId, workload,
			Review.ExpectedGrade, Review.ActualGrade, difficulty, Review.SemesterTaken,
			Review.Lecturer, Review.Content, Review.Suggestion)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}

func GetReviewByModuleAndUser(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		moduleId := strings.ToUpper(c.Param("moduleid"))
		userId, convErr := strconv.Atoi(c.Param("userid"))
		if convErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Userid is malformed",
			})
			return
		}

		exist, review := database.GetReviewByModuleAndUser(db, moduleId, userId)
		if !exist {
			c.JSON(http.StatusOK, gin.H{
				"review": nil,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"review": review,
		})
	}
}

func EditReviewByModuleAndUser(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		moduleId := strings.ToUpper(c.Param("moduleid"))
		userId, convErr := strconv.Atoi(c.Param("userid"))
		if convErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Userid is malformed",
			})
			return
		}

		var EditedReview struct {
			Workload      *int    `json:"workload"`
			ExpectedGrade *string `json:"expected_grade"`
			ActualGrade   *string `json:"actual_grade"`
			Difficulty    *int    `json:"difficulty"`
			SemesterTaken *string `json:"semester_taken"`
			Lecturer      *string `json:"lecturer"`
			Content       *string `json:"content"`
			Suggestion    *string `json:"suggestion"`
		}

		bodyErr := c.ShouldBindJSON(&EditedReview)

		if bodyErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Request body is malformed",
			})
			return
		}

		err := database.EditReviewByModuleAndUser(db, moduleId, userId, EditedReview.Workload,
			EditedReview.ExpectedGrade, EditedReview.ActualGrade, EditedReview.Difficulty,
			EditedReview.SemesterTaken, EditedReview.Lecturer, EditedReview.Content,
			EditedReview.Suggestion)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}

func DeleteReviewByModuleAndUser(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		moduleId := strings.ToUpper(c.Param("moduleid"))
		userId, convErr := strconv.Atoi(c.Param("userid"))
		if convErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"cause":  "Userid is malformed",
			})
			return
		}

		err := database.DeleteReviewByModuleAndUser(db, moduleId, userId)

		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}

func GetWorkload(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"workload": []int{1, 2, 3, 4, 5},
		})
	}
}

func GetGrade(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"grade": []string{"N/A", "A+", "A", "A-", "B+", "B", "B-", "C+", "C", "D+", "D", "F", "S", "U", "CS", "CU", "EXE", "IC", "IP", "W"},
		})
	}
}

func GetDifficulty(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"difficulty": []int{1, 2, 3, 4, 5},
		})
	}
}

func GetSemester(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"semester": []string{"AY2023/2024 S1", "AY2022/2023 S2", "AY2022/2023 S1", "AY2021/2022 S2", "AY2021/2022 S1", "AY2020/2021 S2", "AY2020/2021 S1"},
		})
	}
}

func GetReviewByUser(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		userId, convErr := strconv.Atoi(c.Param("userid"))
		if convErr != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failure",
				"cause":  convErr.Error(),
			})
			return
		}
		review := database.GetReviewByUser(db, userId)
		c.JSON(http.StatusOK, gin.H{
			"review": review,
		})
	}
}
