package vld

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"reflect"
)

type Validator struct {
	validator *validator.Validate
}

func New() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) Validate(i interface{}) map[string]string {
	if err := v.validator.Struct(i); err != nil {
		var validationErrors validator.ValidationErrors
		errorMessages := make(map[string]string)
		if errors.As(err, &validationErrors) {
			for _, e := range validationErrors {
				field := getJSONFieldName(i, e.StructField())
				errorMessages[field] = e.Tag()
			}
			return errorMessages
		}
		errorMessages["error"] = err.Error()
		return errorMessages
	}
	return nil
}

func getJSONFieldName(i interface{}, fieldName string) string {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if field, ok := t.FieldByName(fieldName); ok {
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			return jsonTag
		}
	}

	return fieldName
}
