package config

// Must automatically panics if the error is not nil.
func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}

// MustUnmarshal is a utility function that panics if the unmarshal operation fails.
func MustUnmarshal[T any](unmarshal func([]byte, any) error, src []byte) T {
	var value T

	err := unmarshal(src, &value)
	if err != nil {
		panic(err)
	}

	return value
}
