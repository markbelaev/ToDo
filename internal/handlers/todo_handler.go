package handlers

import (
	"GIN/internal/database"
	"GIN/internal/repository"
	"log/slog"
	"net/http"

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
