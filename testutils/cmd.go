package testutils

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

type CMDResult struct {
	Err     error
	Success bool

	STDOut string
	STDErr string
}

type CMDConfig struct {
	CmdFn  func(t *testing.T)
	MainFn func(t *testing.T, res *CMDResult)
	Env    []string
}

// RunCMD setups a pattern to run test in controlled environments.
// https://stackoverflow.com/a/33404435/9021186
//
// IMPORTANT: this must be run only once per test.
func RunCMD(t *testing.T, config *CMDConfig) {
	t.Helper()

	if os.Getenv("JUST_CHECKING") == "bruh" {
		config.CmdFn(t)
		// Ensure proper exit anyway, otherwise we create a memory leak.
		os.Exit(0)
		return
	}

	// Filter out reserved env.
	env := append(append(os.Environ(), config.Env...), "JUST_CHECKING=bruh")

	outWriter, outCapture, err := CreateSTDCapture(t)
	require.NoError(t, err)
	errWriter, errCapture, err := CreateSTDCapture(t)
	require.NoError(t, err)

	cmd := exec.Command(os.Args[0], "-test.run="+t.Name()) //nolint:gosec
	cmd.Stdout = outWriter
	cmd.Stderr = errWriter
	cmd.Env = env

	err = cmd.Run()

	report := &CMDResult{
		Err:     err,
		Success: true,
		STDOut:  outCapture(),
		STDErr:  errCapture(),
	}

	if err != nil {
		var exErr *exec.ExitError
		if errors.As(err, &exErr) {
			report.Success = exErr.Success()
		} else {
			report.Success = false
		}
	}

	config.MainFn(t, report)
}
