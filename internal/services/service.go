package services

import (
	"log/slog"

	"github.com/fentezi/translator/internal/repositories"
	"github.com/haguro/elevenlabs-go"
)

type Service struct {
	RedisRepository      repositories.Repository
	PostgreSQLRepository *repositories.PostgreSQLRepository
	ClientLabs           *elevenlabs.Client
	log                  *slog.Logger
}

func NewService(rp repositories.Repository, pr *repositories.PostgreSQLRepository, log *slog.Logger, client *elevenlabs.Client) *Service {
	return &Service{
		ClientLabs:           client,
		RedisRepository:      rp,
		log:                  log,
		PostgreSQLRepository: pr,
	}
}
