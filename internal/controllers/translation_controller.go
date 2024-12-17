package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type AddTranslateRequest struct {
	Word        string `json:"text"`
	Translation string `json:"translation"`
}

type TranslateRequest struct {
	Text string `json:"text" validate:"required"`
}

type TranslateResponse struct {
	Translation string `json:"translation"`
}

func (h *Controller) AddTranslate(c echo.Context) error {
	var req AddTranslateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	err := h.service.AddTranslation(req.Word, req.Translation)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusCreated)

}

func (h *Controller) TranslateWord(c echo.Context) error {
	var req TranslateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Field 'text' is required"})
	}

	translation, err := h.service.GetTranslation(req.Text)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch translation"})
	}

	return c.JSON(http.StatusOK, TranslateResponse{Translation: translation})
}
