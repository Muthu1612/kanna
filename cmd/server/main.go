package main

import (
	"log/slog"
	"os"

	"github.com/Muthu1612/kanna/internal/config"
	"github.com/Muthu1612/kanna/internal/logger"
	"github.com/Muthu1612/kanna/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New()

	app := server.New(cfg, log)

	if err := app.Run(); err != nil {
		log.Error("server failed", slog.Any("error", err))
		os.Exit(1)
	}
}
