package testutils

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func CreateSTDCapture(t *testing.T) (writer *os.File, capture func() string, err error) {
	t.Helper()

	reader, writer, err := os.Pipe()
	if err != nil {
		return
	}

	outC := make(chan string)

	// Copy the output in a separate goroutine, so printing can't block indefinitely.
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, reader)
		outC <- buf.String()
	}()

	capture = func() string {
		_ = writer.Close()
		return <-outC
	}

	return
}
