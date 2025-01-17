package main

import (
	"context"
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

	log := logger.NewLogger(cfg.Env)

	log.Info("config", slog.Any("cfg", cfg))

	ctx := context.Background()

	postgres, err := repositories.NewPostgreSQL(cfg)
	if err != nil {
		panic(err)
	}

	postgresDB := repositories.NewPostgreSQLRepository(postgres, ctx)
	log.Info("PostgresDB repository initialized")

	redis, err := repositories.NewRedis(cfg)
	if err != nil {
		panic(err)
	}

	redisDB := repositories.NewRedisRepository(redis, ctx)
	log.Info("RedisDB repository initialized")

	clientLabs := elevenlabs.NewElevenLabs(ctx, cfg.ApiKey)
	log.Info("ElevenLabs client initialized")

	service := services.NewService(redisDB, postgresDB, log, clientLabs)
	log.Info("Service initialized")

	controllers := controllers.NewControllers(service)
	log.Info("Controllers initialized")

	server := server.NewServer(*controllers)
	log.Info("Server initialized")

	e := server.Start(log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Error("shutting down the server")
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown:", slog.Any("err", err))
	} else {
		log.Info("Server shut down gracefully.")
	}

}
