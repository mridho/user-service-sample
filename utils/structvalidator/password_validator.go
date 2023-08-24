package structvalidator

import (
	"user-service-sample/utils/password"

	"github.com/go-playground/validator/v10"
)

const (
	// custom validation tags (prefix with "_")
	passwordTag = "_password"

	passwordValidationMessage = "{0} must contain at least 1 capital characters, 1 number, and 1 special (non alpha-numeric) character"
)

func validatePassword(fl validator.FieldLevel) bool {
	return password.Validate(fl.Field().String())
}

func WithPasswordValidationTag() func(*StructValidator) {

	return func(sv *StructValidator) {
		sv.validator.RegisterValidation(passwordTag, validatePassword)
		ct := customTranslation{
			Tag:         passwordTag,
			Translation: passwordValidationMessage,
		}
		ct.RegisterCustomTranslation(sv.validator, sv.translator)
	}
}
