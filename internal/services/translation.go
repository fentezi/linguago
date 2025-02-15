package services

import (
	"errors"
	"fmt"
	"os"

	"log/slog"

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
	logger := s.log.With(
		slog.String("word", word),
		slog.String("translation", translation),
		slog.String("word_id", wordID.String()),
	)

	logger.Debug("adding translation")

	err := s.PostgreSQLRepository.Set(wordID, word, translation)
	if err != nil {
		logger.Error("failed to save translation to PostgreSQL", slog.Any("error", err))
		return nil, err
	}

	res := &models.Word{
		ID:          wordID,
		Word:        word,
		Translation: translation,
	}

	if err := s.SaveAudio(word, wordID.String()); err != nil {
		logger.Warn("translation saved, but failed to generate audio", slog.Any("error", err))
		return res, nil
	}

	logger.Info("translation added successfully")
	return res, nil
}

func (s *Service) GetTranslation(word string) (string, error) {
	logger := s.log.With(slog.String("word", word))
	logger.Debug("fetching translation")

	translation, err := google.TranslateWordAPI(word)
	if err != nil {
		logger.Error("failed to fetch translation from Google API", slog.Any("error", err))
		return "", fmt.Errorf("failed to get translation from Google API: %w", err)
	}

	wordID := uuid.New()
	saveErr := s.PostgreSQLRepository.Set(wordID, word, translation)
	if saveErr != nil {
		logger.Warn("translation fetched but failed to save to PostgreSQL", slog.Any("error", saveErr))
	}

	if err := s.SaveAudio(word, wordID.String()); err != nil {
		logger.Warn("translation saved but failed to generate audio", slog.Any("error", err))
	}

	logger.Info("translation fetched and cached", slog.String("word_id", wordID.String()))
	return translation, nil
}

func (s *Service) SaveAudio(word, wordID string) error {
	logger := s.log.With(
		slog.String("word", word),
		slog.String("word_id", wordID),
	)

	logger.Debug("generating audio")

	filePath := fmt.Sprintf("./audio/%s.mp3", wordID)
	if _, err := os.Stat(filePath); err == nil {
		logger.Debug("audio file already exists, skipping generation")
		return nil
	} else if !os.IsNotExist(err) {
		logger.Error("failed to check if audio file exists", slog.Any("error", err))
		return fmt.Errorf("failed to check audio file: %w", err)
	}

	if err := elevenlabs.TextToSpeech(s.ClientLabs, wordID, word); err != nil {
		logger.Error("failed to generate audio", slog.Any("error", err))
		return fmt.Errorf("failed to save audio from text: %w", err)
	}

	logger.Info("audio generated successfully", slog.String("file_path", filePath))
	return nil
}

func (s *Service) GetAudio(wordID string) (*os.File, error) {
	logger := s.log.With(slog.String("word_id", wordID))
	logger.Debug("fetching audio file")

	filePath := fmt.Sprintf("./audio/%s.mp3", wordID)
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			logger.Warn("audio file not found")
			return nil, ErrAudioNotFound
		}
		logger.Error("error checking audio file", slog.Any("error", err))
		return nil, fmt.Errorf("failed to read audio file: %w", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		logger.Error("failed to open audio file", slog.Any("error", err))
		return nil, fmt.Errorf("failed to open audio file: %w", err)
	}

	logger.Info("audio file opened", slog.String("file_path", filePath))
	return file, nil
}
