package chans

import (
	"time"
)

type Waiter[T any] struct {
	src       <-chan T
	condition func(T) bool
	timeout   time.Duration
	onClose   func()

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
		res:       make(chan T, 1), // Channel will only hold one result.
	}

	return waiter.init()
}

func (waiter *Waiter[T]) listen() {
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
	if waiter.onClose != nil {
		defer waiter.onClose()
	}

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
