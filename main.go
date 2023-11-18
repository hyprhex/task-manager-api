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

	// Update a specific task
	router.PUT("/tasks/:id", func(ctx *gin.Context) {

		id := ctx.Param("id")

		var updatedTask Task

		if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for i, task := range tasks {
			if task.ID == id {
				if updatedTask.Title != "" {
					tasks[i].Title = updatedTask.Title
				}
				if updatedTask.Description != "" {
					tasks[i].Description = updatedTask.Description
				}

				ctx.JSON(http.StatusOK, gin.H{"message": "Task updated"})
				return
			}
		}

		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})

	})

	// Delete a specific task 
	router.DELETE("/tasks/:id", func(ctx *gin.Context) {

		id := ctx.Param("id")

		for i, val := range tasks {
			if val.ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				ctx.JSON(http.StatusOK, gin.H{"message": "Task removed"})
				return
			}
		}

		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})

	})

	// Create Task
	router.POST("/tasks", func(ctx *gin.Context) {

		var newTask Task

		if err := ctx.ShouldBindJSON(&newTask); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tasks = append(tasks, newTask)
		ctx.JSON(http.StatusCreated, gin.H{"message": "Task created"})

	})

	router.Run() // Listen and serve on 0.0.0.0:8080
}
