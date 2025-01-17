package responses

import "github.com/google/uuid"

type TranslationResponse struct {
	ID          uuid.UUID `json:"id"`
	Word        string    `json:"word"`
	Translation string    `json:"translation"`
}

