package main

import (
	"time"
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
	{ID: "2", Title: "Task 2", Description: "Second Task", DueDate: time.Now().AddDate(0, 0, 1), Status: "Pending"},
	{ID: "3", Title: "Task 3", Description: "Third Task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Pending"},
}

func main() {
	router := gin.Default()
	router.GET("/ping", func (ctx *gin.Context)  {
		ctx.JSON(200, gin.H{
			"message": "poing",
		})	
	})

	router.Run() // Listen and serve on 0.0.0.0:8080
}
