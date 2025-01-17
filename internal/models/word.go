package models

import "github.com/google/uuid"

type Word struct {
	ID          uuid.UUID
	Word        string
	Translation string
}
