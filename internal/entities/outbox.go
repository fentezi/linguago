package entities

import "github.com/google/uuid"

type OutboxMessage struct {
	EventID     uuid.UUID
	Word        string
	Translation string
}
