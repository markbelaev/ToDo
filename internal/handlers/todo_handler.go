package handlers

import (
	"GIN/internal/database"
	"GIN/internal/models"
	"GIN/internal/repository"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ToDoHandler struct {
	repo *repository.ToDoRepository
}

func NewToDoHandler() *ToDoHandler {
	pool := database.GetPool()

	return &ToDoHandler{
		repo: repository.NewToDoRepository(pool),
	}
}

func (h *ToDoHandler) GetToDos(c *gin.Context) {
	ctx := c.Request.Context()

	todos, err := h.repo.GetToDos(ctx)
	if err != nil {
		slog.Error("GetToDos err: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"todos": todos,
	})
}

func (h *ToDoHandler) GetToDoByID(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		slog.Error("GetToDoByID error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	todo, err := h.repo.GetToDoByID(ctx, id)
	if err != nil {
		slog.Error("GetToDoByID error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error getting the task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"todo": todo,
	})
}

func (h *ToDoHandler) DeleteToDo(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		slog.Error("DeleteToDo error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	err = h.repo.DeleteToDo(ctx, id)
	if err != nil {
		slog.Error("DeleteToDo error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error deleting the task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"todo": "deleted",
	})
}

func (h *ToDoHandler) CreateToDo(c *gin.Context) {
	ctx := c.Request.Context()

	var todo models.ToDo
	if err := c.ShouldBindJSON(&todo); err != nil {
		slog.Error("CreateToDo error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	createdToDo, err := h.repo.CreateToDo(ctx, &todo)
	if err != nil {
		slog.Error("CreateToDo error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creating task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"todo created": createdToDo,
	})
}

func (h *ToDoHandler) UpdateToDo(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		slog.Error("UpdateToDo error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	var todo models.ToDo
	if err := c.ShouldBindJSON(&todo); err != nil {
		slog.Error("UpdateToDo error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	todo.ID = int(id)

	updatedToDo, err := h.repo.PutToDo(ctx, &todo)
	if err != nil {
		slog.Error("UpdateToDo error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error updating the task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"todo updated": updatedToDo,
	})
}
