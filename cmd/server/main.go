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
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      bool      `json:"status"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	slog.Info("Connecting to database")

	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file", "error", err)
		os.Exit(1)
	}

	connectString := os.Getenv("DB_CONNECTION_STRING")
	if connectString == "" {
		slog.Error("Error connecting to database")
		os.Exit(1)
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
		slog.Info("GET /todos")

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		rows, err := pool.Query(ctx, "SELECT id, title, description, status, created_at, updated_at FROM todos")
		if err != nil {
			slog.Error("failed to execute query", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get rows",
			})
			return
		}
		defer rows.Close()

		var todos []ToDo
		for rows.Next() {
			var todo ToDo
			if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.Created_at, &todo.Updated_at); err != nil {
				slog.Error("failed to scan row", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to get rows",
				})
				return
			}

			todos = append(todos, todo)
		}

		slog.Info("Successfully fetched all todos")
		c.JSON(http.StatusOK, gin.H{
			"todos": todos,
		})
	}
}

func GetToDoIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("GET /todos/:id")

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		paramID := c.Param("id")
		id, err := strconv.Atoi(paramID)
		if err != nil {
			slog.Error("failed to parse param id", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid id",
			})
			return
		}

		var todo ToDo
		err = pool.QueryRow(ctx, "SELECT id, title, description, status, created_at, updated_at FROM todos WHERE id = $1", id).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.Created_at, &todo.Updated_at)
		if err != nil {
			slog.Error("failed to execute query", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "todo not found",
			})
			return
		}

		slog.Info("Successfully /todos/:id")
		c.JSON(http.StatusOK, gin.H{
			"todo": todo,
		})
	}
}

func PostToDoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("POST /todos")

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		var todo ToDo
		if err := c.BindJSON(&todo); err != nil {
			slog.Error("failed to bind request", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "failed to bind json",
			})
			return
		}

		var id int
		err := pool.QueryRow(ctx, "INSERT INTO todos (title, description, status) VALUES ($1, $2, $3) RETURNING id", todo.Title, todo.Description, todo.Status).Scan(&id)
		if err != nil {
			slog.Error("failed to execute query", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to insert row",
			})
			return
		}

		slog.Info("Successfully /todos/:id")
		c.JSON(http.StatusCreated, gin.H{
			"todo created": todo,
		})
	}
}

func DeleteToDoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("DELETE /todos")

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		paramID := c.Param("id")
		id, err := strconv.Atoi(paramID)
		if err != nil {
			slog.Error("failed to parse param id", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid id",
			})
			return
		}

		_, err = pool.Exec(ctx, "DELETE FROM todos WHERE id = $1", id)
		if err != nil {
			slog.Error("failed to execute query", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to delete row",
			})
			return
		}

		slog.Info("Successfully /todos/:id")
		c.JSON(http.StatusOK, gin.H{
			"todo": "deleted",
		})
	}
}

func PutToDoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("PUT /todos")

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		paramID := c.Param("id")
		id, err := strconv.Atoi(paramID)
		if err != nil {
			slog.Error("failed to parse param id", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid id",
			})
			return
		}

		var todo ToDo
		if err := c.BindJSON(&todo); err != nil {
			slog.Error("failed to bind request", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to bind json",
			})
			return
		}

		_, err = pool.Exec(ctx, "UPDATE todos SET title = $1, description = $2, status = $3, created_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = $4", todo.Title, todo.Description, todo.Status, id)
		if err != nil {
			slog.Error("failed to execute query", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to update row",
			})
			return
		}

		slog.Info("Successfully /todos/:id")
		c.JSON(http.StatusOK, gin.H{
			"todo updated": todo,
		})
	}
}
