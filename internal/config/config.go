package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	LLM      LLMConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type LLMConfig struct {
	Provider string
	Ollama   OllamaConfig
}

type OllamaConfig struct {
	URL   string
	Model string
}

type DatabaseConfig struct {
	URL string
}

func Load() (Config, error) {
	// Load .env for local development.
	// Existing environment variables take precedence.
	_ = godotenv.Load()

	cfg := Config{
		Server: ServerConfig{
			Host: getEnv("KANNA_HOST", "0.0.0.0"),
			Port: getEnv("KANNA_PORT", "8080"),
		},
		LLM: LLMConfig{
			Provider: getEnv("KANNA_LLM_PROVIDER", "ollama"),
			Ollama: OllamaConfig{
				URL:   getEnv("KANNA_OLLAMA_URL", "http://localhost:11434"),
				Model: getEnv("KANNA_OLLAMA_MODEL", ""),
			},
		},
		Database: DatabaseConfig{
			URL: getEnv("KANNA_DATABASE_URL", ""),
		},
	}

	if cfg.LLM.Ollama.Model == "" {
		return Config{}, fmt.Errorf("KANNA_OLLAMA_MODEL is required")
	}

	if cfg.Database.URL == "" {
		return Config{}, fmt.Errorf("KANNA_DATABASE_URL is required")
	}

	if cfg.Server.Host == "" {
		return Config{}, fmt.Errorf("KANNA_HOST is required")
	}

	if cfg.Server.Port == "" {
		return Config{}, fmt.Errorf("KANNA_PORT is required")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}

	return fallback
}
