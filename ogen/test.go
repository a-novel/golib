package ogen

import (
	"context"
	"fmt"
	"net"
)

// MustGetResponse is a helper function to assert the response-type from an ogen-generated client method.
// It returns an error if the response is not of the expected type.
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

// GetRandomPort returns an available random port.
func GetRandomPort() (int, error) {
	listener, err := (&net.ListenConfig{}).Listen(context.Background(), "tcp", ":0")
	if err != nil {
		return 0, fmt.Errorf("failed to create listener: %w", err)
	}

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return 0, fmt.Errorf("expected TCPAddr, got %T", listener.Addr())
	}

	port := addr.Port

	err = listener.Close()
	if err != nil {
		return 0, fmt.Errorf("failed to close listener: %w", err)
	}

	return port, nil
}
