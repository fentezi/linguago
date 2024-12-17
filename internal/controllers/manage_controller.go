package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Controller) GetAllWords(c echo.Context) error {
	words, err := h.service.GetAllWords()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, words)
}

func (h *Controller) DeleteWord(c echo.Context) error {
	word := c.Param("word")

	err := h.service.DeleteTranslation(word)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
