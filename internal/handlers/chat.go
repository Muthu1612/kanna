package handlers

import (
	"net/http"

	"github.com/Muthu1612/kanna/internal/conversation"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	conversationService conversation.Service
}

func NewChatHandler(
	conversationService conversation.Service,
) *ChatHandler {
	return &ChatHandler{
		conversationService: conversationService,
	}
}

type chatRequest struct {
	Message string `json:"message" binding:"required"`
}

type chatResponse struct {
	Message string `json:"message"`
}

func (h *ChatHandler) Chat(c *gin.Context) {
	var request chatRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "message is required",
		})
		return
	}

	response, err := h.conversationService.Chat(
		c.Request.Context(),
		request.Message,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to process conversation",
		})
		return
	}

	c.JSON(http.StatusOK, chatResponse{
		Message: response,
	})
}
