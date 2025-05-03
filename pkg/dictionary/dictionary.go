package dictionary

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type Dictionary struct {
	Phonetic string `json:"phonetic"`
}

const URL = "https://api.dictionaryapi.dev/api/v2/entries/en/"

func Phonetic(word string) (string, error) {
	baseURL, err := url.Parse(URL)
	if err != nil {
		return "", err
	}

	baseURL.Path += url.PathEscape(word)

	resp, err := http.Get(baseURL.String())
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var d []Dictionary
	if err := json.Unmarshal(body, &d); err != nil {
		return "", err
	}

	return d[0].Phonetic, nil

}
