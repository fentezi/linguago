package controllers

import (
	"github.com/fentezi/translator/pkg/helper"
	"net/http"

	"github.com/fentezi/translator/internal/controllers/requests"
	"github.com/labstack/echo/v4"
)

func (h *Controller) GetWords(c echo.Context) error {
	words, err := h.service.GetWords()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, helper.NewAPIError(err))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"words": words,
	})
}

func (h *Controller) DeleteWord(c echo.Context) error {
	var req requests.DeleteWordRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helper.InvalidParam(err))
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, helper.InvalidRequestData(errs))
	}
	err := h.service.DeleteWord(req.WordID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.NewAPIError(err))
	}

	return c.NoContent(http.StatusNoContent)
}
