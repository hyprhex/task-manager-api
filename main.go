package main

import (
	"context"
	"fmt"
	"database/sql"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
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

var db *bun.DB

func main() {

	// Connect to the PostgreSQL database
	ctx := context.Background()
	dsn := "postgres://postgres:root@localhost:5432/task_manager_api?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db = bun.NewDB(sqldb, pgdialect.New())

	// Create the tasks table if not exist 
	_, err := db.NewCreateTable().Model((*Task)(nil)).IfNotExists().Exec(ctx)

	if err != nil {
		fmt.Println("Failed to create table:", err)
		return
	}

	// Add a query for logging 
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	// Ping the database to test the connections 

	err = db.Ping()

	if err != nil {
		fmt.Println("Failed to connect to the database")
		return
	}
	
	// Connection successful
	fmt.Println("Connected to database")

	router := gin.Default()

	// Route handler
	router.GET("/", home)
	router.GET("/tasks", getTasks)
	router.GET("/tasks/:id", getTask)
	router.PUT("/tasks/:id", updateTask)
	router.DELETE("/tasks/:id", deleteTask)
	router.POST("/tasks", createTask)
	
	router.Run() // Listen and serve on 0.0.0.0:8080

}

// Home page 
func home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to the task manager API"})
}

// Getting All tasks 
func getTasks(ctx *gin.Context) {

	var tasks []Task

	err := db.NewSelect().Model(&tasks).Scan(ctx.Request.Context())

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})

}

// Get a specific task 
func getTask(ctx *gin.Context) {

	taskID :=	ctx.Param("id")

	if taskID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID must be present"})
		return
	}

	task := &Task{}

	err := db.NewSelect().Model(task).Where("id = ?", taskID).Scan(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if task.ID == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	ctx.JSON(http.StatusOK, task)



}

// Update a specific task
func updateTask(ctx *gin.Context) {

	taskID := ctx.Param("id")

	if taskID == "" {
		ctx.JSON(http.StatusNoContent, gin.H{"error": "ID must be present"})
		return
	}

	updatedTask := &Task{}

	if err := ctx.ShouldBindJSON(updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.NewUpdate().Model(updatedTask).
		Set("title = ?", updatedTask.Title).
		Set("description = ?", updatedTask.Description).
		Where("id = ?", taskID).
		Exec(ctx.Request.Context())

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task updated"})

}

	// Delete a specific task 
func deleteTask(ctx *gin.Context) {

	taskID := ctx.Param("id")

	task := &Task{}

	res, err := db.NewDelete().Model(task).Where("id = ?", taskID).Exec(ctx.Request.Context())

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected > 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "Task removed"})
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
	}

}

	// Create Task
func createTask(ctx *gin.Context) {

	newTask := &Task{}

	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err :=	db.NewInsert().Model(newTask).Exec(ctx.Request.Context())

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Task created"})


}

