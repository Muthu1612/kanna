package api

import (
	"github.com/gin-gonic/gin"

	"github.com/Muthu1612/kanna/internal/handlers"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/health", handlers.Health)

	return r
}
