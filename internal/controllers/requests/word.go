package requests

import "github.com/google/uuid"

type DeleteWordRequest struct {
	WordID uuid.UUID `param:"word_id" validate:"uuid"`
}

type GetAudioRequest struct {
	WordID uuid.UUID `param:"word_id" validate:"uuid"`
}
