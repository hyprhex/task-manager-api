package main

import (
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
)

// Tasks parameter 
type Task struct {
	ID						string `json:"id"` 
	Title					string `json:"title"`
	Description   string `json:"description"`
	DueDate				time.Time `json:"due_date"`
	Status				string `json:"status"`
}

// Simple tasks 
var tasks = []Task {
	{ID: "1", Title: "Task 1", Description: "First Task", DueDate: time.Now(), Status: "Pending"},
	{ID: "2", Title: "Task 2", Description: "Second Task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
	{ID: "3", Title: "Task 3", Description: "Third Task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
}

func main() {
	router := gin.Default()
	
	// Getting All tasks 
	router.GET("/tasks", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
	})

	// Get a specific task 
	router.GET("/tasks/:id", func(ctx *gin.Context) {

		id := ctx.Param("id")

		for _, task := range tasks {
			if task.ID == id {
				ctx.JSON(http.StatusOK, task)
				return
			}
		}

		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})

	})

	router.Run() // Listen and serve on 0.0.0.0:8080
}
