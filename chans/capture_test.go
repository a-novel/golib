package chans_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/chans"
)

func TestCapture(t *testing.T) {
	t.Parallel()

	source := make(chan int)
	defer close(source)

	capture := chans.NewCapture(source)

	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()

		source <- 1
	}()

	go func() {
		defer wg.Done()

		source <- 2
	}()

	go func() {
		defer wg.Done()

		source <- 3
	}()

	wg.Wait()

	require.Eventually(t, func() bool {
		values := capture.GetAll()

		return assert.Len(t, values, 3) &&
			assert.Contains(t, values, 1) &&
			assert.Contains(t, values, 2) &&
			assert.Contains(t, values, 3)
	}, 100*time.Millisecond, 10*time.Millisecond, "Capture did not receive all values")
}
