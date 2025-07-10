package config

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func LoadEnv[T any](value string, fallback T, parser func(string) (T, error)) T {
	if value == "" {
		return fallback
	}

	parsedValue, err := parser(value)
	if err != nil {
		return fallback
	}

	return parsedValue
}

func SliceParser[T any](parser func(string) (T, error)) func(string) ([]T, error) {
	return func(value string) ([]T, error) {
		// To force an empty slice if the value is empty.
		if value == "[]" {
			return nil, nil
		}

		// Split source on commas.
		parts := strings.Split(value, ",")

		// Parse each part using the provided parser.
		parsedValues := make([]T, 0, len(parts))

		for i, part := range parts {
			trimmedPart := strings.TrimSpace(part)
			if trimmedPart == "" {
				continue // Skip empty parts.
			}

			parsedValue, err := parser(trimmedPart)
			if err != nil {
				// If any parsing fails, deem the whole slice invalid and return the fallback.
				return nil, err
			}

			parsedValues[i] = parsedValue
		}

		// If no parts were parsed, return the fallback.
		if len(parsedValues) == 0 {
			return nil, fmt.Errorf(`value "%s" is empty`, value)
		}

		return parsedValues, nil
	}
}

func StringParser(value string) (string, error) {
	return value, nil
}

func Int64Parser(value string) (int64, error) {
	return strconv.ParseInt(value, 0, 64)
}

func Int32Parser(value string) (int32, error) {
	parsedValue, err := strconv.ParseInt(value, 0, 32)
	if err != nil {
		return 0, err
	}

	return int32(parsedValue), nil
}

func Int16Parser(value string) (int16, error) {
	parsedValue, err := strconv.ParseInt(value, 0, 16)
	if err != nil {
		return 0, err
	}

	return int16(parsedValue), nil
}

func Int8Parser(value string) (int8, error) {
	parsedValue, err := strconv.ParseInt(value, 0, 8)
	if err != nil {
		return 0, err
	}

	return int8(parsedValue), nil
}

func IntParser(value string) (int, error) {
	return strconv.Atoi(value)
}

func Uint64Parser(value string) (uint64, error) {
	parsedValue, err := strconv.ParseUint(value, 0, 64)
	if err != nil {
		return 0, err
	}

	return parsedValue, nil
}

func Uint32Parser(value string) (uint32, error) {
	parsedValue, err := strconv.ParseUint(value, 0, 32)
	if err != nil {
		return 0, err
	}

	return uint32(parsedValue), nil
}

func Uint16Parser(value string) (uint16, error) {
	parsedValue, err := strconv.ParseUint(value, 0, 16)
	if err != nil {
		return 0, err
	}

	return uint16(parsedValue), nil
}

func Uint8Parser(value string) (uint8, error) {
	parsedValue, err := strconv.ParseUint(value, 0, 8)
	if err != nil {
		return 0, err
	}

	return uint8(parsedValue), nil
}

func UintParser(value string) (uint, error) {
	parsedValue, err := strconv.ParseUint(value, 0, 0)
	if err != nil {
		return 0, err
	}

	return uint(parsedValue), nil
}

func BoolParser(value string) (bool, error) {
	parsedValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, err
	}

	return parsedValue, nil
}

func Float64Parser(value string) (float64, error) {
	parsedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return parsedValue, nil
}

func Float32Parser(value string) (float32, error) {
	parsedValue, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return 0, err
	}

	return float32(parsedValue), nil
}

func DurationParser(value string) (time.Duration, error) {
	parsedValue, err := time.ParseDuration(value)
	if err != nil {
		return 0, err
	}

	return parsedValue, nil
}

func TimeParser(value string) (time.Time, error) {
	parsedValue, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, err
	}

	return parsedValue, nil
}

func JSONMapParser(value string) (map[string]any, error) {
	var parsedValue map[string]any

	err := json.Unmarshal([]byte(value), &parsedValue)
	if err != nil {
		return nil, err
	}

	return parsedValue, nil
}

func JSONSliceParser(value string) ([]any, error) {
	var parsedValue []any

	err := json.Unmarshal([]byte(value), &parsedValue)
	if err != nil {
		return nil, err
	}

	return parsedValue, nil
}
