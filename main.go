package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Define a struct for a to-do item
type ToDoItem struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
}

// Define a struct for the application
// (For now, just a wrapper around the database connection)
type App struct {
	DB *sql.DB
}

func main() {
	// Connect to the PostgreSQL database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", "5432", "fsattari", "null", "todolist") // hard coded values for my machine :)
	fmt.Println("Connecting to database:", dsn)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Check if the connection is successful
	if err := db.Ping(); err != nil {
		panic(err)
	}

	// Create the todos table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS todos (id SERIAL PRIMARY KEY, task TEXT NOT NULL)")
	if err != nil {
		panic(err)
	}

	// Create instance of the App struct
	app := &App{DB: db}

	// Initialize the Gin router
	r := gin.Default()

	// Endpoint to add a new to-do item
	r.POST("/todos", app.newTodo)

	// Endpoint to view all to-do items
	r.GET("/todos", app.getTodos)

	// Endpoint to delete a to-do item by ID
	r.DELETE("/todos/:id", app.deleteTodo)

	// Start the server
	r.Run(":8080")
}

// Handler to get all to-do items
func (app *App) getTodos(c *gin.Context) {
	rows, err := app.DB.Query("SELECT id, task FROM todos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}
	defer rows.Close()

	var todos []ToDoItem
	for rows.Next() {
		var item ToDoItem
		if err := rows.Scan(&item.ID, &item.Task); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse items"})
			return
		}
		todos = append(todos, item)
	}

	c.JSON(http.StatusOK, todos)
}

// Handler to add a new to-do item
func (app *App) newTodo(c *gin.Context) {
	var newItem ToDoItem
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	_, err := app.DB.Exec("INSERT INTO todos (task) VALUES ($1)", newItem.Task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Item added"})
}

// Handler to delete a to-do item by ID
func (app *App) deleteTodo(c *gin.Context) {
	id := c.Param("id")

	// Check if the item exists
	var exists bool
	err := app.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM todos WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check item existence"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Proceed to delete the item
	_, err = app.DB.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item successfully deleted"})
}
