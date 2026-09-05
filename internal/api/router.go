package api

import (
	"github.com/Muthu1612/kanna/internal/app"
	"github.com/Muthu1612/kanna/internal/handlers"
	"github.com/gin-gonic/gin"
)

func NewRouter(application *app.App) *gin.Engine {
	r := gin.Default()

	r.GET("/health", handlers.Health)

	chatHandler := handlers.NewChatHandler(application.LLM)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/chat", chatHandler.Chat)
	}

	return r
}
