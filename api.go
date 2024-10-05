package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func api() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	api := r.Group("/api")
	api.GET("/test", func(c *gin.Context) {
		// name := c.Query("name")
		// value := c.Query("value")
		name1 := c.Param("name")
		fmt.Printf("name1: ", name1)
		// fmt.Printf("name: %v\n", name)
		// fmt.Printf("value: %v\n", value)
	})

	api.POST("/register", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		test := c.PostForm("test")

		fmt.Printf("username: %v\n", username)
		fmt.Printf("password: %v\n", password)
		fmt.Printf("test: %v\n", test)
	})

	api.POST("/login", func(c *gin.Context) {
		type LoginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var req LoginRequest
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(200, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, req)

	})

	r.Run(":8080")
}
