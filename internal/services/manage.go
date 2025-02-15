package services

import (
	"fmt"
	"os"

	"log/slog"

	"github.com/fentezi/translator/internal/models"
)

func (s *Service) GetAllWords() ([]models.Word, error) {
	logger := s.log.With(slog.String("operation", "GetAllWords"))
	logger.Debug("fetching all words from PostgreSQL")

	words, err := s.PostgreSQLRepository.Gets()
	if err != nil {
		logger.Error("failed to fetch words from PostgreSQL", slog.Any("error", err))
		return nil, err
	}

	logger.Info("fetched words from PostgreSQL", slog.Int("count", len(words)))
	return words, nil
}

func (s *Service) DeleteTranslation(wordID string) error {
	logger := s.log.With(slog.String("word_id", wordID), slog.String("operation", "DeleteTranslation"))
	logger.Debug("starting deletion process")

	if err := s.PostgreSQLRepository.Delete(wordID); err != nil {
		logger.Error("failed to delete translation from PostgreSQL", slog.Any("error", err))
		return err
	}
	logger.Debug("deleted translation from PostgreSQL")

	filePath := fmt.Sprintf("./audio/%s.mp3", wordID)
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Warn("audio file not found, skipping deletion", slog.String("file_path", filePath))
		} else {
			logger.Error("failed to delete audio file", slog.String("file_path", filePath), slog.Any("error", err))
			return fmt.Errorf("failed to delete audio file: %w", err)
		}
	} else {
		logger.Debug("deleted audio file", slog.String("file_path", filePath))
	}

	logger.Info("translation deleted successfully")
	return nil
}
