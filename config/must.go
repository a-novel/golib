package config

func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}

func MustUnmarshal[T any](unmarshal func([]byte, any) error, src []byte) T {
	var value T

	err := unmarshal(src, &value)
	if err != nil {
		panic(err)
	}

	return value
}
