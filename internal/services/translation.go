package services

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/fentezi/translator/pkg/elevenlabs"
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

	s.log.Debug("translation added successfully", slog.String("word", word), slog.String("translation", translation))
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

func (s *Service) GetAudio(word string) (*os.File, error) {
	s.log.Debug("request to generate audio started", slog.String("word", word))
	filePath := fmt.Sprintf("./audio/%s.mp3", word)

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			err := elevenlabs.TextToSpeech(s.ClientLabs, word)
			if err != nil {
				s.log.Error("failed to generate audio from text", slog.String("word", word), slog.Any("error", err))
				return nil, fmt.Errorf("failed to generate audio from text: %w", err)
			}
			s.log.Debug("audio generated successfully", slog.String("word", word))

		}
	}

	s.log.Debug("attempting to open audio file", slog.String("filePath", filePath))

	file, err := os.Open(filePath)
	if err != nil {
		s.log.Error("failed to open audio file", slog.String("word", word), slog.String("filePath", filePath), slog.Any("error", err))
		return nil, fmt.Errorf("failed to open audio file: %w", err)
	}

	s.log.Debug("audio file opened successfully", slog.String("word", word), slog.String("filePath", filePath))
	return file, nil
}
