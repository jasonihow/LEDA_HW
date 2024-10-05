package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func api_hw() {
	r := gin.Default()
	todo := r.Group("/todo")

	todo.POST("/add", func(c *gin.Context) {

		var req TODO
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"收到請求": req})
		TODOlist = append(TODOlist, req)
	})
	todo.GET("/list", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"list": TODOlist})
	})
	todo.DELETE("/delete/:task", func(c *gin.Context) {
		task := c.Param("task")
		for i, todo := range TODOlist {
			if todo.Task == task {
				TODOlist = append(TODOlist[:i], TODOlist[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"成功刪除task: ": task})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"找不到task: ": task})
	})

	todo.PUT("/update/:task", func(c *gin.Context) {
		task := c.Param("task")
		status := c.Query("status")

		var done bool
		if status == "true" {
			done = true
		} else if status == "false" {
			done = false
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "status must be true or false"})
			return
		}

		var req TODO
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for i, todo := range TODOlist {
			if todo.Task == task {
				TODOlist[i].Done = done
				c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully", "task": req})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
	})

	todo.GET("/complete", func(c *gin.Context) {
		var completedTasks []TODO
		for _, todo := range TODOlist {
			if todo.Done {
				completedTasks = append(completedTasks, todo)
			}
		}
		c.JSON(http.StatusOK, gin.H{"已完成事項": completedTasks})
	})

	todo.GET("/uncomplete", func(c *gin.Context) {
		var uncompletedTasks []TODO
		for _, todo := range TODOlist {
			if !todo.Done {
				uncompletedTasks = append(uncompletedTasks, todo)
			}
		}
		c.JSON(http.StatusOK, gin.H{"待辦事項": uncompletedTasks})
	})

	todo.GET("/clear", func(c *gin.Context) {
		TODOlist = []TODO{}
		c.JSON(http.StatusOK, gin.H{"message": "已成功清除所有task"})
	})

	r.Run(":8080")

}

var TODOlist = []TODO{}

type TODO struct {
	Task string `json:"task"`
	Done bool   `json:"done"`
}
