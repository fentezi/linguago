package google

import (
	gtranslate "github.com/gilang-as/google-translate"
)

// translate text from English to Russian
func TranslateWordAPI(text string) (string, error) {
	value := gtranslate.Translate{
		Text: text,
		From: "en",
		To:   "ru",
	}

	translated, err := gtranslate.Translator(value)
	if err != nil {
		return "", err
	}

	return translated.Text, nil

}
