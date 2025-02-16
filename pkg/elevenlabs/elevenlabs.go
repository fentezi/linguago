package elevenlabs

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/haguro/elevenlabs-go"
)

const (
	ModelID = "eleven_turbo_v2"
	VoiceID = "9BWtsMINqrJLrRacOk9x"
)

func New(ctx context.Context, apiKey string) *elevenlabs.Client {
	client := elevenlabs.NewClient(ctx, apiKey, 10*time.Second)

	return client
}

func TextToSpeech(client *elevenlabs.Client, wordID uuid.UUID, word string) error {
	ttsReq := elevenlabs.TextToSpeechRequest{
		Text:    word,
		ModelID: ModelID,
		VoiceSettings: &elevenlabs.VoiceSettings{
			SpeakerBoost:    false,
			SimilarityBoost: 0.75,
			Stability:       1.0,
		},
	}

	audio, err := client.TextToSpeech(
		VoiceID,
		ttsReq,
	)

	if err != nil {
		return err
	}

	if err := os.WriteFile(fmt.Sprintf("./audio/%s.mp3", wordID), audio, 0644); err != nil {
		return err
	}

	return nil
}
