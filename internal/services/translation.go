package services

import (
	"fmt"
	"log/slog"

	google "github.com/fentezi/translator/pkg/google_translate"
)

func (s *Service) AddTranslation(word, translation string) error {
	err := s.PostgreSQLRepository.Set(word, translation)
	if err != nil {
		s.log.Error("failed to save translation to PostgreSQL", slog.String("word", word), slog.Any("error", err))
		return err
	}

	err = s.RedisRepository.Set(word, translation)
	if err != nil {
		s.log.Error("failed to cache translation in Redis", slog.String("word", word), slog.Any("error", err))
		return err
	}

	s.log.Debug("translation added successfully", slog.String("word", word))
	return nil
}

func (s *Service) GetTranslation(word string) (string, error) {
	translation, err := s.RedisRepository.Get(word)
	if err == nil {
		s.log.Debug("cache hit", slog.String("word", word))
		return translation, nil
	}
	s.log.Debug("cache miss", slog.String("word", word))

	translation, err = google.TranslateWordAPI(word)
	if err != nil {
		s.log.Error("failed to fetch translation from Google API", slog.String("word", word), slog.Any("error", err))
		return "", fmt.Errorf("failed to get translation from Google API: %w", err)
	}

	saveErr := s.PostgreSQLRepository.Set(word, translation)
	if saveErr != nil {
		s.log.Error("failed to save translation to PostgreSQL", slog.String("word", word), slog.Any("error", saveErr))
		return "", fmt.Errorf("failed to save translation to PostgreSQL: %w", saveErr)
	}

	cacheErr := s.RedisRepository.Set(word, translation)
	if cacheErr != nil {
		s.log.Error("failed to cache translation in Redis", slog.String("word", word), slog.Any("error", cacheErr))
		return "", fmt.Errorf("failed to cache translation in Redis: %w", cacheErr)
	}

	s.log.Debug("translation fetched and cached successfully", slog.String("word", word))
	return translation, nil
}
