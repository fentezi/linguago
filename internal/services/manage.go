package services

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/fentezi/translator/internal/models"
)

func (s *Service) GetAllWords() ([]models.Word, error) {
	s.log.Debug("fetching all words from PostgreSQL")

	words, err := s.PostgreSQLRepository.Gets()
	if err != nil {
		s.log.Error("failed to fetch all words from PostgreSQL", slog.Any("error", err))
		return nil, err
	}

	s.log.Debug("successfully fetched all words from PostgreSQL", slog.Int("count", len(words)))
	return words, nil
}

func (s *Service) DeleteTranslation(wordID string) error {
	s.log.Debug("attempting to delete translation", slog.String("word_id", wordID))

	err := s.RedisRepository.Delete(wordID)
	if err != nil {
		s.log.Error("failed to delete translation from Redis", slog.String("word_id", wordID), slog.Any("error", err))
		return err
	}
	s.log.Debug("successfully deleted translation from Redis", slog.String("word", wordID))

	err = s.PostgreSQLRepository.Delete(wordID)
	if err != nil {
		s.log.Error("failed to delete translation from PostgreSQL", slog.String("word_id", wordID), slog.Any("error", err))
		return err
	}
	s.log.Debug("successfully deleted translation from PostgreSQL", slog.String("wordID", wordID))

	filePath := fmt.Sprintf("./audio/%s.mp3", wordID)
	err = os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		s.log.Error("failed to delete audio file", slog.String("word_id", wordID), slog.String("filePath", filePath), slog.Any("error", err))
		return fmt.Errorf("failed to delete audio file: %w", err)
	}

	s.log.Debug("translation deleted successfully", slog.String("wordIDD", wordID))
	return nil
}
