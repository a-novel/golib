package database

import (
	"reflect"

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

func registerSortDirectionType(field reflect.Value) interface{} {
	value, ok := field.Interface().(SortDirection)
	if !ok {
		return nil
	}

	return string(value)
}

// RegisterSortDirection registers the SortDirection type with a validator.
func RegisterSortDirection(validator *validator.Validate) {
	validator.RegisterCustomTypeFunc(registerSortDirectionType, SortDirection(""))
}
