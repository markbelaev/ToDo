package routes

import (
	"GIN/internal/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	todoHandler := handlers.NewToDoHandler()

	api := router.Group("/api")
	{
		api.GET("/todos", todoHandler.GetToDos)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
}
