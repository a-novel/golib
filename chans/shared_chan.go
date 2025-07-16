package chans

import (
	"sync"
)

// Shared forwards the data of a single channel to multiple listeners.
type Shared[T any] struct {
	// The listeners that are registered to receive data from the source channel.
	listeners map[chan T]struct{}

	mu sync.RWMutex
}

func NewShared[T any]() *Shared[T] {
	multi := &Shared[T]{listeners: make(map[chan T]struct{})}

	return multi
}

func (multi *Shared[T]) readMSG(msg T) {
	// Go channels are thread-safe. However, a channel can be closed
	// while a message is being sent to it.
	// Thr read-lock prevents the closing of listeners while a message is being sent.
	multi.mu.RLock()
	defer multi.mu.RUnlock()

	// Forward the message to all listeners.
	for listener := range multi.listeners {
		listener <- msg
	}
}

// Send a new message to the source channel. This will be forwarded to all listeners.
func (multi *Shared[T]) Send(data T) {
	multi.readMSG(data)
}

// Register a new listener, and return it. The listener will receive all messages sent to the source channel.
func (multi *Shared[T]) Register() <-chan T {
	multi.mu.Lock()
	defer multi.mu.Unlock()

	listener := make(chan T)
	multi.listeners[listener] = struct{}{} // mark active.

	return listener
}

// Unregister a listener. The listener will no longer receive messages from the source channel.
func (multi *Shared[T]) Unregister(src <-chan T) {
	multi.mu.Lock()
	defer multi.mu.Unlock()

	for listener := range multi.listeners {
		if listener == src {
			delete(multi.listeners, listener)
			close(listener)
		}
	}
}

// Close the source channel, along with all listeners.
func (multi *Shared[T]) Close() {
	multi.mu.Lock()
	defer multi.mu.Unlock()

	for listener := range multi.listeners {
		delete(multi.listeners, listener)
		close(listener)
	}
}
