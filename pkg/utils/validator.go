package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(req interface{}) map[string]string {
	err := validate.Struct(req)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		field := strings.ToLower(err.Field())
		errors[field] = err.Tag()
	}
	return errors
}
