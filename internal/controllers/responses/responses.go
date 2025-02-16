package responses

import "github.com/google/uuid"

type CreateWordResponse struct {
	ID          uuid.UUID `json:"id"`
	Word        string    `json:"word"`
	Translation string    `json:"translation"`
}
