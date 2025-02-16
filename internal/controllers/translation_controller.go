package controllers

import (
	"errors"
	"net/http"

	"github.com/fentezi/translator/internal/controllers/requests"
	"github.com/fentezi/translator/internal/controllers/responses"
	"github.com/fentezi/translator/internal/repositories"
	"github.com/labstack/echo/v4"
)

func (h *Controller) CreateWord(c echo.Context) error {
	var req requests.CreateWordRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := h.service.CreateWord(req.Word, req.Translation)
	if err != nil {
		if errors.Is(err, repositories.ErrAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, "word already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	translation, err := h.service.TranslateWord(req.Word)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch translation")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"translation": translation,
	})
}

func (h *Controller) GetAudio(c echo.Context) error {
	var req requests.GetAudioRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	file, err := h.service.GetAudio(req.WordID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch audio")
	}

	defer file.Close()

	return c.Stream(http.StatusOK, "audio/mp3", file)
}
