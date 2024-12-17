package services

import (
	"log/slog"

	"github.com/fentezi/translator/internal/repositories"
)

type Service struct {
	RedisRepository      repositories.Repository
	PostgreSQLRepository *repositories.PostgreSQLRepository
	log                  *slog.Logger
}

func NewService(rp repositories.Repository, pr *repositories.PostgreSQLRepository, log *slog.Logger) *Service {
	return &Service{
		RedisRepository:      rp,
		log:                  log,
		PostgreSQLRepository: pr,
	}
}
