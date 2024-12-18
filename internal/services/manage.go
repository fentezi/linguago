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

func (s *Service) DeleteTranslation(word string) error {
	s.log.Debug("attempting to delete translation", slog.String("word", word))

	err := s.RedisRepository.Delete(word)
	if err != nil {
		s.log.Error("failed to delete translation from Redis", slog.String("word", word), slog.Any("error", err))
		return err
	}
	s.log.Debug("successfully deleted translation from Redis", slog.String("word", word))

	err = s.PostgreSQLRepository.Delete(word)
	if err != nil {
		s.log.Error("failed to delete translation from PostgreSQL", slog.String("word", word), slog.Any("error", err))
		return err
	}
	s.log.Debug("successfully deleted translation from PostgreSQL", slog.String("word", word))

	filePath := fmt.Sprintf("./audio/%s.mp3", word)
	err = os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		s.log.Error("failed to delete audio file", slog.String("word", word), slog.String("filePath", filePath), slog.Any("error", err))
		return fmt.Errorf("failed to delete audio file: %w", err)
	}

	s.log.Debug("translation deleted successfully", slog.String("word", word))
	return nil
}
