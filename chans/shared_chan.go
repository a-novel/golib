package chans

import (
	"sync"
)

// Shared forwards the data of a single channel to multiple listeners.
type Shared[T any] struct {
	// The source channel
	src chan T
	// The listeners that are registered to receive data from the source channel. The boolean key is used to indicate
	// whether the listener is still active or not.
	listeners map[chan T]bool

	mu sync.Mutex
}

func NewShared[T any]() *Shared[T] {
	multi := &Shared[T]{src: make(chan T), listeners: make(map[chan T]bool)}
	go multi.listen()

	return multi
}

func (multi *Shared[T]) readMsg(msg T) {
	multi.mu.Lock()
	defer multi.mu.Unlock()

	for listener, ok := range multi.listeners {
		// Listener is not active anymore.
		if !ok {
			continue
		}

		listener <- msg
	}
}

func (multi *Shared[T]) listen() {
	// Forward each new message to all listeners.
	for msg := range multi.src {
		multi.readMsg(msg)
	}
}

// Chan returns the source channel of the Shared.
func (multi *Shared[T]) Chan() chan<- T {
	return multi.src
}

// Send a new message to the source channel. This will be forwarded to all listeners.
func (multi *Shared[T]) Send(data T) {
	multi.src <- data
}

// Register a new listener, and return it. The listener will receive all messages sent to the source channel.
func (multi *Shared[T]) Register() <-chan T {
	multi.mu.Lock()
	defer multi.mu.Unlock()

	listener := make(chan T)
	multi.listeners[listener] = true // mark active.

	return listener
}

// Unregister a listener. The listener will no longer receive messages from the source channel.
func (multi *Shared[T]) Unregister(src <-chan T) {
	multi.mu.Lock()
	defer multi.mu.Unlock()

	for listener, ok := range multi.listeners {
		if listener == src && ok {
			delete(multi.listeners, listener)
			close(listener)
		}
	}
}

// Close the source channel, along with all listeners.
func (multi *Shared[T]) Close() {
	close(multi.src)

	multi.mu.Lock()
	defer multi.mu.Unlock()

	for listener, ok := range multi.listeners {
		if ok {
			delete(multi.listeners, listener)
			close(listener)
		}
	}
}
