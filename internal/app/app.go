package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Muthu1612/kanna/internal/config"
	"github.com/Muthu1612/kanna/internal/conversation"
	"github.com/Muthu1612/kanna/internal/llm"
	"github.com/Muthu1612/kanna/internal/memory"
)

type App struct {
	LLM                 llm.Client
	Memory              memory.Store
	ConversationService conversation.Service
}

func New(cfg config.Config, logger *slog.Logger) (*App, error) {
	var llmClient llm.Client

	switch cfg.LLM.Provider {
	case "ollama":
		llmClient = llm.NewOllamaClient(
			cfg.LLM.Ollama.URL,
			cfg.LLM.Ollama.Model,
		)

	default:
		return nil, fmt.Errorf(
			"unsupported LLM provider: %s",
			cfg.LLM.Provider,
		)
	}

	ctx := context.Background()

	memoryStore, err := memory.NewPostgresStore(
		ctx,
		cfg.Database.URL,
	)
	if err != nil {
		return nil, fmt.Errorf("initialize memory: %w", err)
	}

	conversationService := conversation.NewService(
		llmClient,
		memoryStore,
	)

	return &App{
		LLM:                 llmClient,
		Memory:              memoryStore,
		ConversationService: conversationService,
	}, nil
}
