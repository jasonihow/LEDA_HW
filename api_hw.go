package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Config struct for configuration
type Config struct {
	DBName string `json:"db_name"`
	DBUser string `json:"db_user"`
	DBPass string `json:"db_pass"`
}

// TODO struct
type TODO struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}

// Global variables
var (
	db      *gorm.DB
	todoMux sync.RWMutex
)

// LoadConfig loads configuration from config.json
func LoadConfig() (Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

// Initialize the database and create default tasks
func InitializeDatabase(config Config) {
	var err error
	db, err = gorm.Open(sqlite.Open(config.DBName), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto migrate
	db.AutoMigrate(&TODO{})

	// Create default tasks if the table is empty
	var count int64
	db.Model(&TODO{}).Count(&count)
	if count == 0 {
		defaultTasks := []TODO{
			{Task: "Sample Task 1", Done: false},
			{Task: "Sample Task 2", Done: false},
		}
		db.Create(&defaultTasks)
	}
}

func api_hw() {
	// Load the configuration
	config, err := LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Initialize the database
	InitializeDatabase(config)

	r := gin.Default()
	todo := r.Group("/api/todo")

	todo.POST("/", createTODO)
	todo.GET("/", listTODOs)
	todo.DELETE("/:id", deleteTODO)
	todo.PUT("/:id", updateTODO)
	todo.GET("/complete", listCompletedTODOs)
	todo.GET("/uncomplete", listUncompletedTODOs)
	todo.DELETE("/clear", clearTODOs)

	r.Run(":8008")
}

// Create new TODO
func createTODO(c *gin.Context) {
	var req TODO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": 4001, "message": err.Error()})
		return
	}

	todoMux.Lock()
	defer todoMux.Unlock()

	db.Create(&req)
	c.JSON(http.StatusCreated, gin.H{"message": "Task created successfully", "task": req})
}

// List all TODOs
func listTODOs(c *gin.Context) {
	var todos []TODO

	todoMux.RLock()
	defer todoMux.RUnlock()

	db.Find(&todos)
	c.JSON(http.StatusOK, gin.H{"todos": todos})
}

// Delete a TODO
func deleteTODO(c *gin.Context) {
	id := c.Param("id")

	todoMux.Lock()
	defer todoMux.Unlock()

	var todo TODO
	if result := db.First(&todo, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error_code": 4041, "message": "Task not found"})
		return
	}

	// Delete the specified TODO item
	db.Delete(&todo)

	// Renumber the remaining TODO items
	var todos []TODO
	if err := db.Order("id").Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve remaining tasks"})
		return
	}
	for i, t := range todos {
		// Update the ID directly using the address of the item
		newID := uint(i + 1) // New ID
		t.ID = newID         // Update the ID

		// Save the updated TODO item
		if err := db.Model(&t).Update("ID", newID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to renumber tasks"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// Update specified TODO
func updateTODO(c *gin.Context) {
	id := c.Param("id")
	var req TODO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": 4002, "message": err.Error()})
		return
	}

	todoMux.Lock()
	defer todoMux.Unlock()

	var todo TODO
	if result := db.First(&todo, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error_code": 4042, "message": "Task not found"})
		return
	}

	todo.Task = req.Task
	todo.Done = req.Done

	db.Save(&todo)
	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully", "task": todo})
}

// List completed TODOs
func listCompletedTODOs(c *gin.Context) {
	var completedTodos []TODO

	todoMux.RLock()
	defer todoMux.RUnlock()

	db.Where("done = ?", true).Find(&completedTodos)
	c.JSON(http.StatusOK, gin.H{"completed_todos": completedTodos})
}

// List uncompleted TODOs
func listUncompletedTODOs(c *gin.Context) {
	var uncompletedTodos []TODO

	todoMux.RLock()
	defer todoMux.RUnlock()

	db.Where("done = ?", false).Find(&uncompletedTodos)
	c.JSON(http.StatusOK, gin.H{"uncompleted_todos": uncompletedTodos})
}

// Clear all TODOs
func clearTODOs(c *gin.Context) {
	todoMux.Lock()
	defer todoMux.Unlock()

	db.Exec("DELETE FROM todos") // Clear the table

	// Reset auto-increment ID
	db.Exec("DELETE FROM sqlite_sequence WHERE name='todos'") // For SQLite

	c.JSON(http.StatusOK, gin.H{"message": "All tasks cleared successfully"})
}
