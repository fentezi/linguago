package helper

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
)

type APIError struct {
	Msg any `json:"msg"`
}

func (e *APIError) Error() string {
	return fmt.Sprint(e.Msg)
}

func NewAPIError(err error) *APIError {
	return &APIError{Msg: err.Error()}
}

func InvalidJSON() *APIError {
	return NewAPIError(fmt.Errorf("invalid JSON request data"))
}

func InvalidParam(err error) *APIError {
	var bindingErr *echo.HTTPError
	if errors.As(err, &bindingErr) {
		return &APIError{Msg: bindingErr.Message}
	}
	return NewAPIError(fmt.Errorf("invalid request"))
}

func InvalidRequestData(errs map[string]string) *APIError {
	errorMessages := ""
	for field, err := range errs {
		errorMessages += fmt.Sprintf("%s: %s;", field, err)
	}

	return &APIError{Msg: errorMessages}
}
