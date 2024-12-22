package main

import (
	"context"
	"log/slog"

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
	log.Info("Server started")

	server.Start(log)

}
