package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

/*
type module struct {
	Code             string `json:"code"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Subscriber_Count int    `json:"subscriber_count"`
}

var example_module_list = []module{
	{Code: "CS3233", Name: "Competitive Programming", Description: "This module aims to prepare students in competitive problem solving. It covers techniques for attacking and solving challenging computational problems. Fundamental algorithmic solving techniques covered include divide and conquer, greedy, dynamic programming, backtracking and branch and bound. Domain specific techniques like number theory, computational geometry, string processing and graph theoretic will also be covered. Advanced AI search techniques like iterative deepening, A* and heuristic search will be included. The module also covers algorithmic and programming language toolkits used in problem solving supported by the solution of representative or well-known problems in the various algorithmic paradigms.", Subscriber_Count: 0},
}

func Get_module(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"module_list": example_module_list,
	})
}

*/
