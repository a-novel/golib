package database

import (
	"github.com/go-playground/validator/v10"
)

// SortDirection controls the direction of the ordering for a particular request. You can use this type with a
// validator.
type SortDirection string

const (
	SortDirectionNone SortDirection = ""
	SortDirectionAsc  SortDirection = "asc"
	SortDirectionDesc SortDirection = "desc"
)

// ValidateEnum creates a custom validation for go-validator. It checks if the value is part of the enum.
//
// TODO: look for custom errors in v11: https://github.com/go-playground/validator/issues/669
func ValidateEnum[T comparable](list ...T) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(T)
		if !ok {
			return false
		}

		for _, element := range list {
			if value == element {
				return true
			}
		}

		return false
	}
}

func MustRegisterValidation(
	customValidator *validator.Validate, name string, validationFn func(fl validator.FieldLevel) bool,
) {
	err := customValidator.RegisterValidation(name, validationFn)
	if err != nil {
		panic(err)
	}
}

// RegisterSortDirection registers the SortDirection type with a validator.
func RegisterSortDirection(customValidator *validator.Validate) {
	MustRegisterValidation(
		customValidator,
		"sort_direction",
		ValidateEnum(SortDirectionNone, SortDirectionAsc, SortDirectionDesc),
	)
}
