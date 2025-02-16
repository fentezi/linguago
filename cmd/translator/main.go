package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fentezi/translator/config"
	"github.com/fentezi/translator/internal/controllers"
	"github.com/fentezi/translator/internal/repositories"
	"github.com/fentezi/translator/internal/server"
	"github.com/fentezi/translator/internal/services"
	"github.com/fentezi/translator/pkg/elevenlabs"
	"github.com/fentezi/translator/pkg/logger"
)

func main() {
	cfg := config.MustConfig()

	log := logger.NewLogger(cfg.Env).With(slog.String("component", "main"))
	log.Info("starting application", slog.Any("config", cfg))

	ctx := context.Background()

	repo, err := repositories.New(ctx, &cfg.Postgres)
	if err != nil {
		log.Error("failed to initialize repository", slog.Any("error", err))
		panic(err)
	}
	log.Info("repository initialized")
	defer repo.Close()

	clientLabs := elevenlabs.New(ctx, cfg.ApiKey)
	log.Info("ElevenLabs client initialized")

	service := services.New(repo, log, clientLabs)
	log.Info("service layer initialized")

	controller := controllers.New(service)
	log.Info("controller initialized")

	srv := server.New(*controller)
	log.Info("server initialized")

	e := srv.Start(log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		log.Info("starting HTTP server on :8080")
		if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server encountered a critical error", slog.Any("error", err))
			panic(err)
		}
	}()
	<-ctx.Done()

	log.Info("received shutdown signal, shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Error("server forced to shutdown", slog.Any("error", err))
	} else {
		log.Info("server shut down gracefully")
	}
}
