package grpc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/grpc"
	adaptersmocks "github.com/a-novel/golib/loggers/adapters/mocks"
)

type customExecServiceIn struct {
	In1 string
}

type customExecServiceOut struct {
	Out1 string
}

type customExecServiceImpl struct {
	err error
}

func (s *customExecServiceImpl) Exec(_ context.Context, data *customExecServiceIn) (*customExecServiceOut, error) {
	return &customExecServiceOut{Out1: "got: " + data.In1}, s.err
}

func TestServiceWithMetrics(t *testing.T) {
	fooErr := errors.New("uwups")

	testCases := []struct {
		name string

		in    *customExecServiceIn
		inErr error

		expect    *customExecServiceOut
		expectErr error
	}{
		{
			name: "OK",

			in: &customExecServiceIn{In1: "test"},

			expect: &customExecServiceOut{Out1: "got: test"},
		},
		{
			name: "Error",

			in:    &customExecServiceIn{In1: "test"},
			inErr: fooErr,

			expect:    &customExecServiceOut{Out1: "got: test"},
			expectErr: fooErr,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			logger := adaptersmocks.NewMockGRPC(t)
			logger.On("Report", "my-service", testCase.inErr)

			service := &customExecServiceImpl{err: testCase.inErr}
			serviceWrapped := grpc.ServiceWithMetrics("my-service", service, logger)

			out, err := service.Exec(context.Background(), testCase.in)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, out)

			out, err = serviceWrapped.Exec(context.Background(), testCase.in)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, out)

			logger.AssertExpectations(t)
		})
	}
}
