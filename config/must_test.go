package config_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/config"
)

var errPanic = errors.New("panic")

func TestMust(t *testing.T) {
	t.Parallel()

	fn1 := func() (string, error) {
		return "foo", nil
	}

	fn2 := func() (string, error) {
		return "", errPanic
	}

	require.Equal(t, "foo", config.Must(fn1()))
	require.PanicsWithValue(t, errPanic, func() {
		config.Must(fn2())
	})
}

func TestMustUnmarshal(t *testing.T) {
	t.Parallel()

	unmarshalFunc := func(data []byte, v any) error {
		if string(data) == "valid" {
			*v.(*string) = "unmarshaled"

			return nil
		}

		return errPanic
	}

	validData := []byte("valid")
	invalidData := []byte("invalid")

	require.Equal(t, "unmarshaled", config.MustUnmarshal[string](unmarshalFunc, validData))
	require.PanicsWithValue(t, errPanic, func() {
		config.MustUnmarshal[string](unmarshalFunc, invalidData)
	})
}
