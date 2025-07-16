package chans_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/chans"
)

func TestWaiter(t *testing.T) {
	t.Parallel()

	strChan := make(chan string, 3)
	defer close(strChan)

	waiter := chans.NewWaiter(strChan, func(s string) bool {
		return s == "bar"
	}, 100*time.Millisecond, nil)

	strChan <- "foo"

	_, ok := waiter.Wait()
	require.False(t, ok)

	waiter = chans.NewWaiter(strChan, func(s string) bool {
		return s == "bar"
	}, 100*time.Millisecond, nil)

	strChan <- "foo"

	strChan <- "bar"

	res, ok := waiter.Wait()
	require.True(t, ok)
	require.Equal(t, "bar", res)
}

func TestWaiterOnClose(t *testing.T) {
	t.Parallel()

	strChan := make(chan string)
	defer close(strChan)

	var closed bool

	waiter := chans.NewWaiter(strChan, func(s string) bool {
		return s == "bar"
	}, 100*time.Millisecond, func() {
		closed = true
	})

	strChan <- "foo"

	_, ok := waiter.Wait()
	require.False(t, ok)
	require.True(t, closed)
}
