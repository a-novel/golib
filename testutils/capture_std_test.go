package testutils_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/testutils"
)

func TestCaptureSTD(t *testing.T) {
	outWriter, outCapture, err := testutils.CreateSTDCapture(t)
	require.NoError(t, err)

	originalSTDOut := os.Stdout
	os.Stdout = outWriter

	_, err = os.Stdout.WriteString("Hello, World!")
	require.NoError(t, err)

	os.Stdout = originalSTDOut

	captured := outCapture()

	require.Equal(t, "Hello, World!", captured)
}
