package handlers

import (

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"github.com/AIpill/task-manager-api/model"

)

var DB *bun.DB

// Home page 
func Home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to the task manager API"})
}

// Getting All tasks 
func GetTasks(ctx *gin.Context) {

	var tasks []model.Task

	err := DB.NewSelect().Model(&tasks).Scan(ctx.Request.Context())

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})

}

// Get a specific task 
func GetTask(ctx *gin.Context) {

	taskID :=	ctx.Param("id")

	if taskID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID must be present"})
		return
	}

	task := &model.Task{}

	err := DB.NewSelect().Model(task).Where("id = ?", taskID).Scan(ctx.Request.Context())
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
func UpdateTask(ctx *gin.Context) {

	taskID := ctx.Param("id")

	if taskID == "" {
		ctx.JSON(http.StatusNoContent, gin.H{"error": "ID must be present"})
		return
	}

	updatedTask := &model.Task{}

	if err := ctx.ShouldBindJSON(updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := DB.NewUpdate().Model(updatedTask).
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
func DeleteTask(ctx *gin.Context) {

	taskID := ctx.Param("id")

	task := &model.Task{}

	res, err := DB.NewDelete().Model(task).Where("id = ?", taskID).Exec(ctx.Request.Context())

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
func CreateTask(ctx *gin.Context) {

	newTask := &model.Task{}

	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err :=	DB.NewInsert().Model(newTask).Exec(ctx.Request.Context())

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Task created"})


}

