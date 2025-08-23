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
		v1 := api.Group("/v1")
		v1.GET("/todos", todoHandler.GetToDos)
		v1.GET("/todos/:id", todoHandler.GetToDoByID)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
}
