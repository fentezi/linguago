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
	"github.com/fentezi/translator/pkg/closer"
	"github.com/fentezi/translator/pkg/elevenlabs"
	"github.com/fentezi/translator/pkg/logger"
	"github.com/labstack/gommon/log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	c := &closer.Closer{}
	log.Info("config initializing")
	cfg := config.MustConfig()

	log.Info("logger initializing")
	logging := logger.NewLogger(cfg.Env).With(slog.String("component", "main"))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logging.Info("repository initializing")
	repo, err := repositories.New(ctx, cfg.Postgres)
	if err != nil {
		logging.Error("failed to initialize repository", slog.Any("error", err))
		panic(err)
	}
	c.Add(repo.Close)

	logging.Info("broker initializing")
	broker, err := kafka.New(cfg.Broker)
	if err != nil {
		logging.Error("failed to initialize kafka broker", slog.Any("error", err))
		panic(err)
	}
	c.Add(broker.Close)

	outboxProducer := outboxproducer.New(logging, broker, repo)

	go func() {
		if err := cron(ctx, outboxProducer, logging, cfg.Broker.Topic); err != nil {
			logging.Error("failed to start cron", slog.Any("error", err))
		}
	}()

	logging.Info("elevenlabs initializing")
	clientLabs := elevenlabs.New(ctx, cfg.ApiKey)

	logging.Info("service initializing")
	service := services.New(logging, clientLabs, repo)

	logging.Info("controllers initializing")
	controller := controllers.New(service)

	logging.Info("server initializing")
	srv := server.New(controller)

	c.Add(srv.Close)

	go func() {
		err = srv.Start(logging, cfg.Server)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logging.Error("failed to start server", slog.Any("error", err))
			panic(err)
		}
	}()

	<-ctx.Done()

	logging.Info("received shutdown signal, shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := c.Close(shutdownCtx); err != nil {
		logging.Error("failed to close server", slog.Any("error", err))
	} else {
		logging.Info("server shutdown gracefully")
	}

}

func cron(
	ctx context.Context, producer outboxproducer.OutboxProducer, log *slog.Logger, topic string,
) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			break
		case <-ticker.C:
			if err := producer.ProduceMessage(ctx, topic); err != nil {
				log.Error("failed to produce message", slog.Any("error", err))
			}
		}
	}
}
