package requests

type CreateWordRequest struct {
	Word        string `json:"word" validate:"required"`
	Translation string `json:"translation" validate:"required"`
}

type TranslateWordRequest struct {
	Word string `json:"word" validate:"required"`
}
