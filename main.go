package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ToDo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
}

func main() {
	connStr := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	ctx := context.Background()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("failed to create connection: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping connection: %v", err)
	}

	log.Println("Connected to PostgreSQL")

	r := gin.Default()

	r.GET("/todos", GetToDoHandler(pool))

	r.GET("/todos/:id", GetToDoIDHandler(pool))

	r.POST("/todos", PostToDoHandler(pool))

	r.DELETE("/todos/:id", DeleteToDoHandler(pool))

	r.PUT("/todos/:id", PutToDoHandler(pool))

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func GetToDoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		rows, err := pool.Query(ctx, "SELECT id, name, description, status FROM todos")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get rows",
			})
			return
		}

		var todos []ToDo
		for rows.Next() {
			var todo ToDo
			if err := rows.Scan(&todo.ID, &todo.Name, &todo.Description, &todo.Status); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to get rows",
				})
				return
			}

			todos = append(todos, todo)
		}

		c.JSON(http.StatusOK, gin.H{
			"todos": todos,
		})
	}
}

func GetToDoIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		paramID := c.Param("id")
		id, err := strconv.Atoi(paramID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid id",
			})
			return
		}

		var todo ToDo
		err = pool.QueryRow(ctx, "SELECT id, name, description, status FROM todos WHERE id = $1", id).Scan(&todo.ID, &todo.Name, &todo.Description, &todo.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get rows",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"todo": todo,
		})
	}
}

func PostToDoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		var todo ToDo
		if err := c.BindJSON(&todo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to bind json",
			})
			return
		}

		var id int
		err := pool.QueryRow(ctx, "INSERT INTO todos (name, description, status) VALUES ($1, $2, $3) RETURNING id", todo.Name, todo.Description, todo.Status).Scan(&id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to insert row",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"todo created": todo,
		})
	}
}

func DeleteToDoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		paramID := c.Param("id")
		id, err := strconv.Atoi(paramID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid id",
			})
			return
		}

		_, err = pool.Exec(ctx, "DELETE FROM todos WHERE id = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to delete row",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"todo": "deleted",
		})
	}
}

func PutToDoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		paramID := c.Param("id")
		id, err := strconv.Atoi(paramID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid id",
			})
			return
		}

		var todo ToDo
		if err := c.BindJSON(&todo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to bind json",
			})
			return
		}

		_, err = pool.Exec(ctx, "UPDATE todos SET name = $1, description = $2, status = $3 WHERE id = $4", todo.Name, todo.Description, todo.Status, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to update row",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"todo updated": todo,
		})
	}
}
