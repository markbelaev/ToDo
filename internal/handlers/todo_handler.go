package handlers

import (
	"GIN/internal/database"
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

	todos, err := h.repo.GetAllToDos(ctx)
	if err != nil {
		slog.Error("GetAllToDos error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error getting toDos",
		})
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
