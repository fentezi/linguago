package controllers

import (
	"errors"
	"net/http"

	"github.com/fentezi/translator/internal/repositories"
	"github.com/fentezi/translator/internal/requests"
	"github.com/fentezi/translator/internal/responses"
	"github.com/fentezi/translator/internal/services"
	"github.com/labstack/echo/v4"
)

func (h *Controller) AddTranslate(c echo.Context) error {
	var req requests.AddRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	res, err := h.service.AddTranslation(req.Word, req.Translation)
	if err != nil {
		if errors.Is(err, repositories.ErrAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, "word already exists in list")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := responses.TranslationResponse{
		ID:          res.ID,
		Word:        res.Word,
		Translation: res.Translation,
	}

	return c.JSON(http.StatusCreated, resp)

}

func (h *Controller) TranslateWord(c echo.Context) error {
	var req requests.TranslateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	translation, err := h.service.GetTranslation(req.Word)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch translation")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"translation": translation,
	})
}

func (h *Controller) GetAudio(c echo.Context) error {
	wordID := c.Param("word_id")
	file, err := h.service.GetAudio(wordID)
	if err != nil {
		if errors.Is(err, services.ErrAudioNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, services.ErrAudioNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch audio")
	}

	defer file.Close()

	return c.Stream(http.StatusOK, "audio/mp3", file)
}
	