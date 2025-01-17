package requests

type AddRequest struct {
	Word        string `json:"word" validate:"required"`
	Translation string `json:"translation" validate:"required"`
}

type TranslateRequest struct {
	Word string `json:"word" validate:"required"`
}
