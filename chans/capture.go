package chans

import (
	"crypto/rand"
	"sync"
)

// Capture accumulates messages sent to a channel.
// Those messages can be accessed any moment as a
// slice.
type Capture[T any] struct {
	values sync.Map
	source <-chan T
}

func NewCapture[T any](source <-chan T) *Capture[T] {
	capture := &Capture[T]{source: source}

	go capture.listen()

	return capture
}

func (capture *Capture[T]) listen() {
	// Listen on the main channel. Every time
	// a value is received, it is stored in a
	// thread-safe map.
	// The random key is used to prevent
	// collisions, since std only provides
	// a thread-safe map implementation.
	for value := range capture.source {
		capture.values.Store(rand.Text(), value)
	}
}

func (capture *Capture[T]) GetAll() []T {
	var values []T

	capture.values.Range(func(_, value any) bool {
		tValue, ok := value.(T)
		if !ok {
			// If the value is not of type T, we skip it.
			// This should not happen, but it makes the linters happy.
			return true
		}

		values = append(values, tValue)

		return true
	})

	return values
}
