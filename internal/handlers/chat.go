package handlers

import (
	"net/http"

	"github.com/Muthu1612/kanna/internal/llm"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	llm llm.Client
}

func NewChatHandler(llmClient llm.Client) *ChatHandler {
	return &ChatHandler{
		llm: llmClient,
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

	response, err := h.llm.Generate(
		c.Request.Context(),
		llm.Request{
			Prompt: request.Message,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate response",
		})
		return
	}

	c.JSON(http.StatusOK, chatResponse{
		Message: response.Content,
	})
}
