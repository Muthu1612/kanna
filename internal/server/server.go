package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Muthu1612/kanna/internal/api"
	"github.com/Muthu1612/kanna/internal/config"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func New(cfg config.Config, logger *slog.Logger) *Server {
	router := api.NewRouter()

	return &Server{
		httpServer: &http.Server{
			Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
			Handler: router,
		},
		logger: logger,
	}
}

func (s *Server) Run() error {
	go func() {
		s.logger.Info(
			"Kanna server starting",
			slog.String("address", s.httpServer.Addr),
		)

		if err := s.httpServer.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			s.logger.Error(
				"server failed",
				slog.Any("error", err),
			)
		}
	}()

	s.waitForShutdown()

	return nil
}

func (s *Server) waitForShutdown() {
	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	s.logger.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error(
			"server shutdown failed",
			slog.Any("error", err),
		)

		return
	}

	s.logger.Info("Kanna server stopped")
}
