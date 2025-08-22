package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type ToDo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	slog.Info("Connecting to database")

	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file", "error", err)
		os.Exit(1)
	}

	connectString := os.Getenv("DB_CONNECTED_STRING")
	if connectString == "" {
		slog.Error("Error connecting to database")
	}

	config, err := pgxpool.ParseConfig(connectString)
	if err != nil {
		slog.Error("failed to parse config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		slog.Error("failed to create connection", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("failed to ping connection", "error", err)
		os.Exit(1)
	}

	slog.Info("Connected to PostgreSQL")

	r := gin.Default()

	r.GET("/todos", GetToDoHandler(pool))

	r.GET("/todos/:id", GetToDoIDHandler(pool))

	r.POST("/todos", PostToDoHandler(pool))

	r.DELETE("/todos/:id", DeleteToDoHandler(pool))

	r.PUT("/todos/:id", PutToDoHandler(pool))

	if err := r.Run(":8080"); err != nil {
		slog.Error("failed to run server", "error", err)
		os.Exit(1)
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
