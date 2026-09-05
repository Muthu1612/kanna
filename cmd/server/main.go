package main

import (
	"log/slog"
	"os"

	"github.com/Muthu1612/kanna/internal/app"
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

	application, err := app.New(cfg, log)
	if err != nil {
		log.Error("application initialization failed", slog.Any("error", err))
		os.Exit(1)
	}

	srv := server.New(cfg, log, application)

	if err := srv.Run(); err != nil {
		log.Error("server failed", slog.Any("error", err))
		os.Exit(1)
	}
}
