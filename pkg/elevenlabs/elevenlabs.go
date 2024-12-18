package elevenlabs

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/haguro/elevenlabs-go"
)

const (
	ModelID = "eleven_turbo_v2"
	VoiceID = "9BWtsMINqrJLrRacOk9x"
)

func NewElevenLabs(ctx context.Context, apiKey string) *elevenlabs.Client {
	client := elevenlabs.NewClient(ctx, apiKey, 10*time.Second)

	return client
}

func TextToSpeech(client *elevenlabs.Client, text string) error {
	ttsReq := elevenlabs.TextToSpeechRequest{
		Text:    text,
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

	if err := os.WriteFile(fmt.Sprintf("./audio/%s.mp3", text), audio, 0644); err != nil {
		return err
	}

	return nil
}
