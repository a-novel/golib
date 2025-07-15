package chans_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/chans"
)

func TestShared(t *testing.T) {
	t.Parallel()

	// Create a new shared channel for integers.
	shared := chans.NewShared[int]()
	defer shared.Close()

	// Register 2 listeners to the shared channel.
	listener1 := shared.Register()
	require.NotNil(t, listener1)

	listener2 := shared.Register()
	require.NotNil(t, listener2)

	capturer1 := chans.NewCapture(listener1)
	capturer2 := chans.NewCapture(listener2)

	// Send a value to the shared channel.
	go func() {
		shared.Send(1)
	}()

	go func() {
		shared.Send(2)
	}()

	go func() {
		shared.Send(3)
	}()

	require.Eventually(t, func() bool {
		values1, values2 := capturer1.GetAll(), capturer2.GetAll()

		return assert.Len(t, values1, 3) &&
			assert.Len(t, values2, 3) &&
			assert.Contains(t, values1, 1) &&
			assert.Contains(t, values1, 2) &&
			assert.Contains(t, values1, 3) &&
			assert.Contains(t, values2, 1) &&
			assert.Contains(t, values2, 2) &&
			assert.Contains(t, values2, 3)
	}, 100*time.Millisecond, 10*time.Millisecond, "Listeners did not receive the expected values")

	// Close listener 2.
	shared.Unregister(listener2)

	// Send another value to the shared channel.
	shared.Send(4)

	// Check that listener 1 received the new value, but listener 2 did not.
	require.Eventually(t, func() bool {
		values1, values2 := capturer1.GetAll(), capturer2.GetAll()

		return assert.Len(t, values1, 4) &&
			assert.Len(t, values2, 3) &&
			assert.Contains(t, values1, 4) &&
			assert.NotContains(t, values2, 4)
	}, 100*time.Millisecond, 10*time.Millisecond, "Listener 1 did not receive the new value")
}
