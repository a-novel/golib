package database_test

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/database"
)

func TestSortDirection(t *testing.T) {
	t.Run("Validation", func(t *testing.T) {
		testCases := []struct {
			name string

			value database.SortDirection

			expectErr bool
		}{
			{
				name: "OK/Empty",

				value: database.SortDirectionNone,
			},
			{
				name: "OK/Asc",

				value: database.SortDirectionAsc,
			},
			{
				name: "OK/Desc",

				value: database.SortDirectionDesc,
			},
			{
				name: "KO/Invalid",

				value: "invalid",

				expectErr: true,
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				customValidator := validator.New(validator.WithRequiredStructEnabled())
				database.RegisterSortDirection(customValidator)

				toValidate := struct {
					Value database.SortDirection `validate:"omitempty,oneof=asc desc"`
				}{
					Value: testCase.value,
				}

				err := customValidator.Struct(toValidate)

				if testCase.expectErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			})
		}

		t.Run("OtherTypes", func(t *testing.T) {
			customValidator := validator.New(validator.WithRequiredStructEnabled())
			database.RegisterSortDirection(customValidator)

			toValidate := struct {
				Value interface{} `validate:"omitempty,oneof=asc desc"`
			}{
				Value: "asc",
			}

			err := customValidator.Struct(toValidate)
			require.NoError(t, err)
		})
	})
}
