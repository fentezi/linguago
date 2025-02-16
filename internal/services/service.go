package services

import (
	"log/slog"

	"github.com/fentezi/translator/internal/repositories"
	"github.com/haguro/elevenlabs-go"
)

type Service struct {
	PostgreSQLRepository *repositories.PostgreSQLRepository
	ClientLabs           *elevenlabs.Client
	log                  *slog.Logger
}

func New(pr *repositories.PostgreSQLRepository, log *slog.Logger, client *elevenlabs.Client) *Service {
	return &Service{
		ClientLabs:           client,
		log:                  log,
		PostgreSQLRepository: pr,
	}
}
