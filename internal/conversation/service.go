package conversation

import (
	"context"
	"fmt"

	"github.com/Muthu1612/kanna/internal/agent"
	"github.com/Muthu1612/kanna/internal/memory"
)

type Service interface {
	Chat(ctx context.Context, message string) (string, error)
}

type conversationService struct {
	agent  agent.Agent
	memory memory.Store
}

func NewService(
	agentService agent.Agent,
	memoryStore memory.Store,
) Service {
	return &conversationService{
		agent:  agentService,
		memory: memoryStore,
	}
}

func (s *conversationService) Chat(
	ctx context.Context,
	message string,
) (string, error) {
	response, err := s.agent.Run(
		ctx,
		agent.Request{
			Message: message,
		},
	)
	if err != nil {
		return "", fmt.Errorf("process conversation: %w", err)
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
