package chans

import (
	"crypto/rand"
	"sync"
)

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
			return true
		}

		values = append(values, tValue)

		return true
	})

	return values
}
