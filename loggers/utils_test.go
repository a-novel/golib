package loggers_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/a-novel/golib/testutils"
)

type stdCapture struct {
	stdout string
	stderr string
}

func captureSTD(t *testing.T) (func() stdCapture, error) {
	t.Helper()

	originalSTDOut := os.Stdout
	originalSTDErr := os.Stderr

	stdOut, stdOutCapture, err := testutils.CreateSTDCapture(t)
	if err != nil {
		return nil, fmt.Errorf("capture stdout: %w", err)
	}

	stdErr, stdErrCapture, err := testutils.CreateSTDCapture(t)
	if err != nil {
		return nil, fmt.Errorf("capture stderr: %w", err)
	}

	os.Stdout = stdOut
	os.Stderr = stdErr

	return func() stdCapture {
		res := stdCapture{
			stdout: stdOutCapture(),
			stderr: stdErrCapture(),
		}

		os.Stdout = originalSTDOut
		os.Stderr = originalSTDErr

		return res
	}, nil
}
