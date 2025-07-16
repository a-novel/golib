package chans

import (
	"time"
)

// Waiter listens to a channel until a value that satisfies a condition is received.
type Waiter[T any] struct {
	// The channel to listen to.
	src <-chan T
	// The condition that the received value must satisfy to be returned.
	condition func(T) bool
	// The timeout after which the waiter will stop waiting for a value.
	timeout time.Duration
	// The function to call once a value is received, or the timeout is reached.
	onClose func()

	// The source channel is listened to on a separate goroutine, to prevent blocking.
	// This channel is buffered (non-blocking), so the Wait method can safely
	// use it.
	res chan T
}

func NewWaiter[T any](
	src <-chan T,
	condition func(T) bool,
	timeout time.Duration,
	onClose func(),
) *Waiter[T] {
	waiter := &Waiter[T]{
		src:       src,
		condition: condition,
		timeout:   timeout,
		onClose:   onClose,
		res:       make(chan T, 1), // The channel only needs to hold 1 result.
	}

	return waiter.init()
}

func (waiter *Waiter[T]) listen() {
	// Listen on the source channel. If a value is received that satisfies the condition,
	// send it to the output channel. The main loop will wait on it.
	for msg := range waiter.src {
		if waiter.condition(msg) {
			waiter.res <- msg

			return
		}
	}
}

func (waiter *Waiter[T]) init() *Waiter[T] {
	go waiter.listen()

	return waiter
}

func (waiter *Waiter[T]) Wait() (T, bool) {
	// Run the onClose function if it is set.
	if waiter.onClose != nil {
		defer waiter.onClose()
	}

	// Prevent infinite blocking. If the value is not received within the timeout,
	// quit and return a zero value.
	timer := time.NewTimer(waiter.timeout)
	defer timer.Stop()

	for {
		select {
		case msg := <-waiter.res:
			return msg, true
		case <-timer.C:
			var zero T

			return zero, false
		default:
			continue
		}
	}
}
