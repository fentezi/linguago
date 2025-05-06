package services

import (
	"log/slog"

	"github.com/fentezi/translator/internal/repositories"
	"github.com/haguro/elevenlabs-go"
)

type Service struct {
	Repository repositories.PostgreSQLRepository
	ClientLabs elevenlabs.Client
	log        *slog.Logger
}

func New(
	log *slog.Logger, client elevenlabs.Client, pr repositories.PostgreSQLRepository,
) Service {
	return Service{
		ClientLabs: client,
		log:        log,
		Repository: pr,
	}
}
