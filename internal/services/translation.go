package services

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/fentezi/translator/internal/models"
	"github.com/fentezi/translator/pkg/elevenlabs"
	google "github.com/fentezi/translator/pkg/google_translate"
	"github.com/google/uuid"
)

var (
	ErrAudioNotFound = errors.New("audio not found")
)

func (s *Service) AddTranslation(word, translation string) (*models.Word, error) {
	wordID := uuid.New()
	res, err := s.PostgreSQLRepository.Set(wordID, word, translation)
	if err != nil {
		s.log.Error("failed to save translation to PostgreSQL", slog.String("word", word), slog.Any("error", err))
		return nil, err
	}

	err = s.RedisRepository.Set(word, translation)
	if err != nil {
		s.log.Error("failed to cache translation in Redis", slog.String("word", word), slog.Any("error", err))
		return nil, err
	}
	err = s.SaveAudio(word, wordID.String())
	if err != nil {
		s.log.Error("failed to save audio", slog.String("word", word), slog.Any("error", err))
		return nil, err
	}

	s.log.Debug("translation added successfully", slog.String("word", word), slog.String("translation", translation))
	return res, nil
}

func (s *Service) GetTranslation(word string) (string, error) {
	wordID := uuid.New()

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

	_, saveErr := s.PostgreSQLRepository.Set(wordID, word, translation)
	if saveErr != nil {
		s.log.Error("failed to save translation to PostgreSQL", slog.String("word", word), slog.Any("error", saveErr))
		return "", fmt.Errorf("failed to save translation to PostgreSQL: %w", saveErr)
	}

	cacheErr := s.RedisRepository.Set(word, translation)
	if cacheErr != nil {
		s.log.Error("failed to cache translation in Redis", slog.String("word", word), slog.Any("error", cacheErr))
		return "", fmt.Errorf("failed to cache translation in Redis: %w", cacheErr)
	}
	err = s.SaveAudio(word, wordID.String())
	if err != nil {
		s.log.Error("failed to save audio", slog.String("word", word), slog.Any("error", err))
		return "", fmt.Errorf("failed to save audio: %w", err)
	}

	s.log.Debug("translation fetched and cached successfully", slog.String("word", word))
	return translation, nil
}

func (s *Service) SaveAudio(word, wordID string) error {
	s.log.Debug("request to save audio started", slog.String("word_id", wordID))

	filePath := fmt.Sprintf("./audio/%s.mp3", wordID)

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			err := elevenlabs.TextToSpeech(s.ClientLabs, wordID, word)
			if err != nil {
				s.log.Error("failed to save audio from text", slog.String("word", word), slog.Any("error", err))
				return fmt.Errorf("failed to save audio from text: %w", err)
			}
			s.log.Debug("audio generated successfully", slog.String("word", word))

		}
	}
	s.log.Debug("audio saved successfully", slog.String("word_id", wordID))
	return nil
}

func (s *Service) GetAudio(wordID string) (*os.File, error) {
	s.log.Debug("request to get audio", slog.String("word_id", wordID))
	filePath := fmt.Sprintf("./audio/%s.mp3", wordID)

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrAudioNotFound
		}
		return nil, fmt.Errorf("failed to read audio file: %w", err)
	}

	s.log.Debug("attempting to open audio file", slog.String("filePath", filePath))

	file, err := os.Open(filePath)
	if err != nil {
		s.log.Error("failed to open audio file", slog.String("word_id", wordID), slog.String("filePath", filePath), slog.Any("error", err))
		return nil, fmt.Errorf("failed to open audio file: %w", err)
	}

	s.log.Debug("audio file opened successfully", slog.String("word_id", wordID), slog.String("filePath", filePath))
	return file, nil
}
