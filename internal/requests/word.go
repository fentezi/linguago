package requests

type DeleteRequest struct {
	WordID string `param:"word_id" validate:"uuid"`
}

type GetAudioRequest struct {
	WordID string `param:"word_id" validate:"uuid"`
}
