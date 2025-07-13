package chans

import (
	"sync"
	"time"
)

type Waiter[T any] struct {
	src       <-chan T
	condition func(T) bool
	timeout   time.Duration
	onClose   func()

	res      T
	wg       sync.WaitGroup
	timedOut bool
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
	}

	return waiter.init()
}

func (waiter *Waiter[T]) listen() {
	defer waiter.wg.Done()

	for {
		select {
		case <-time.After(waiter.timeout):
			waiter.timedOut = true

			return
		case msg, ok := <-waiter.src:
			if !ok {
				return // Channel closed, exit the goroutine.
			}

			if waiter.condition(msg) {
				waiter.res = msg

				return
			}
		}
	}
}

func (waiter *Waiter[T]) init() *Waiter[T] {
	waiter.wg.Add(1)

	go waiter.listen()

	return waiter
}

func (waiter *Waiter[T]) Wait() (T, bool) {
	if waiter.onClose != nil {
		defer waiter.onClose()
	}

	waiter.wg.Wait()

	if waiter.timedOut {
		var zero T

		return zero, false
	}

	return waiter.res, true
}
