package main

import (
	"context"
	"fmt"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"

	"github.com/AIpill/task-manager-api/handlers"
	"github.com/AIpill/task-manager-api/model"
)

var db *bun.DB

func main() {

	db, err := connectToDatabase()
	if err != nil {

		fmt.Println("Failed to connect to the database:", err)
		return

	}
	defer db.Close()

	err = createTaskTable(db)
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

	handlers.DB = db
	router := gin.Default()

	// Route handler
	router.GET("/", handlers.Home)
	router.GET("/tasks", handlers.GetTasks)
	router.GET("/tasks/:id", handlers.GetTask)
	router.PUT("/tasks/:id", handlers.UpdateTask)
	router.DELETE("/tasks/:id", handlers.DeleteTask)
	router.POST("/tasks", handlers.CreateTask)
	
	router.Run() // Listen and serve on 0.0.0.0:8080

}

func connectToDatabase() (*bun.DB, error) {

	dsn := "postgres://postgres:root@localhost:5432/task_manager_api?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db = bun.NewDB(sqldb, pgdialect.New())
	
	return db, nil

}

func createTaskTable(db *bun.DB) error {

	ctx := context.Background()
	_, err := db.NewCreateTable().Model((*model.Task)(nil)).IfNotExists().Exec(ctx)

	return err

}
