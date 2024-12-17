package models

type Word struct {
	ID          int    `json:"id"`
	Text        string `json:"text"`
	Translation string `json:"translation"`
}
