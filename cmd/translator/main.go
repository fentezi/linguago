package main

import (
	"context"
	"errors"
	"github.com/fentezi/translator/config"
	"github.com/fentezi/translator/internal/controllers"
	outboxproducer "github.com/fentezi/translator/internal/cron/outbox_producer"
	"github.com/fentezi/translator/internal/kafka"
	"github.com/fentezi/translator/internal/repositories"
	"github.com/fentezi/translator/internal/server"
	"github.com/fentezi/translator/internal/services"
	"github.com/fentezi/translator/pkg/elevenlabs"
	"github.com/fentezi/translator/pkg/logger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	cfg := config.MustConfig()
	log := logger.NewLogger(cfg.Env).With(slog.String("component", "main"))
	log.Info("starting application", slog.Any("config", cfg))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	repo, err := repositories.New(ctx, &cfg.Postgres)
	if err != nil {
		log.Error("failed to initialize repository", slog.Any("error", err))
		panic(err)
	}
	defer repo.Close()
	log.Info("repository initialized")

	broker, err := kafka.New(cfg.Broker)
	if err != nil {
		log.Error("failed to initialize kafka broker", slog.Any("error", err))
		panic(err)
	}

	outboxProducer := outboxproducer.New(log, broker, repo)
	go func() {
		if err := cron(ctx, outboxProducer, log, cfg.Broker.Topic); err != nil {
			log.Error("failed to start cron", slog.Any("error", err))
		}
	}()

	clientLabs := elevenlabs.New(ctx, cfg.ApiKey)
	log.Info("ElevenLabs client initialized")

	service := services.New(repo, log, clientLabs)
	log.Info("service layer initialized")

	controller := controllers.New(service)
	log.Info("controller initialized")

	srv := server.New(*controller)
	log.Info("server initialized")

	e := srv.Start(log)

	go func() {
		log.Info("starting HTTP server on :8090")
		if err := e.Start(":8090"); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

	outboxProducer.Close()
}

func cron(
	ctx context.Context, producer *outboxproducer.OutboxProducer, log *slog.Logger, topic string,
) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("shutting down cron job")
			return nil
		case <-ticker.C:
			if err := producer.ProduceMessage(ctx, topic); err != nil {
				log.Error("failed to produce message", slog.Any("error", err))
			}
		}
	}
}
