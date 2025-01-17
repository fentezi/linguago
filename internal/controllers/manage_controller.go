package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Controller) GetAllWords(c echo.Context) error {
	words, err := h.service.GetAllWords()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(words) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "list words empty")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"words": words,
	})
}

func (h *Controller) DeleteWord(c echo.Context) error {
	word := c.Param("word_id")

	err := h.service.DeleteTranslation(word)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
