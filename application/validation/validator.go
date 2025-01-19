package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate

func NewValidator() {
	Validator = validator.New()

	usernameRegex := regexp.MustCompile("^[a-z0-9][a-z0-9]*([._][a-z0-9]+)*$")

	err := Validator.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		return usernameRegex.Match([]byte(fl.Field().String()))
	})

	if err != nil {
		panic("unable to register username validator")
	}
}
