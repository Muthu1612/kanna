package app

import (
	"fmt"
	"log/slog"

	"github.com/Muthu1612/kanna/internal/config"
	"github.com/Muthu1612/kanna/internal/llm"
)

type App struct {
	LLM llm.Client
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

	return &App{
		LLM: llmClient,
	}, nil
}
