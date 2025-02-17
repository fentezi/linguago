package controllers

import (
	"errors"
	"github.com/fentezi/translator/pkg/helper"
	"net/http"

	"github.com/fentezi/translator/internal/controllers/requests"
	"github.com/fentezi/translator/internal/controllers/responses"
	"github.com/fentezi/translator/internal/repositories"
	"github.com/labstack/echo/v4"
)

func (h *Controller) CreateWord(c echo.Context) error {
	var req requests.CreateWordRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helper.InvalidJSON())
	}
	if errs := h.validator.Validate(req); len(errs) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, helper.InvalidRequestData(errs))
	}

	res, err := h.service.CreateWord(req.Word, req.Translation)
	if err != nil {
		if errors.Is(err, repositories.ErrAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, helper.NewAPIError(errors.New("word already exists")))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, helper.NewAPIError(err))
	}

	resp := responses.CreateWordResponse{
		ID:          res.ID,
		Word:        res.Word,
		Translation: res.Translation,
	}

	return c.JSON(http.StatusCreated, resp)

}

func (h *Controller) TranslateWord(c echo.Context) error {
	var req requests.TranslateWordRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helper.InvalidJSON())
	}
	if errs := h.validator.Validate(req); len(errs) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, helper.InvalidRequestData(errs))
	}

	translation, err := h.service.TranslateWord(req.Word)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, helper.NewAPIError(err))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"translation": translation,
	})
}

func (h *Controller) GetAudio(c echo.Context) error {
	var req requests.GetAudioRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helper.InvalidParam(err))
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, helper.InvalidRequestData(errs))
	}

	file, err := h.service.GetAudio(req.WordID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, helper.NewAPIError(err))
	}

	defer file.Close()

	return c.Stream(http.StatusOK, "audio/mp3", file)
}
