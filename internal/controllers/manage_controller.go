package controllers

import (
	"net/http"

	"github.com/fentezi/translator/internal/controllers/requests"
	"github.com/labstack/echo/v4"
)

func (h *Controller) GetWords(c echo.Context) error {
	words, err := h.service.GetWords()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"words": words,
	})
}

func (h *Controller) DeleteWord(c echo.Context) error {
	var req requests.DeleteWordRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	err := h.service.DeleteWord(req.WordID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
