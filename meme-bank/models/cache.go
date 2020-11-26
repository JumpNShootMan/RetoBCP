package models

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

// Validator returns a cached instance of a validator
func Validator() *validator.Validate {
	if validate == nil {
		validate = validator.New()
	}
	return validate
}
