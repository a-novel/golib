package ogen

import "fmt"

func MustGetResponse[Raw any, Want any](res Raw, err error) (Want, error) {
	var zero Want

	if err != nil {
		return zero, err
	}

	out, ok := any(res).(Want)
	if !ok {
		return zero, fmt.Errorf("unexpected response: %v", res)
	}

	return out, nil
}
