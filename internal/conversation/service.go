package conversation

import (
	"context"
	"fmt"

	"github.com/Muthu1612/kanna/internal/llm"
	"github.com/Muthu1612/kanna/internal/memory"
)

type Service interface {
	Chat(ctx context.Context, message string) (string, error)
}

type conversationService struct {
	llm    llm.Client
	memory memory.Store
}

func NewService(
	llmClient llm.Client,
	memoryStore memory.Store,
) Service {
	return &conversationService{
		llm:    llmClient,
		memory: memoryStore,
	}
}

func (s *conversationService) Chat(
	ctx context.Context,
	message string,
) (string, error) {
	response, err := s.llm.Generate(
		ctx,
		llm.Request{
			Prompt: message,
		},
	)
	if err != nil {
		return "", fmt.Errorf("generate response: %w", err)
	}

	err = s.memory.SaveConversation(
		ctx,
		memory.Conversation{
			Question: message,
			Answer:   response.Content,
		},
	)
	if err != nil {
		return "", fmt.Errorf("save conversation: %w", err)
	}

	return response.Content, nil
}
