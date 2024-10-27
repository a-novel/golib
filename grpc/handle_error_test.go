package grpc_test

import (
	"errors"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/a-novel/golib/grpc"
)

type FooErr struct {
	message string
}

func (e *FooErr) Error() string {
	return e.message
}

func NewFooErr(message string) error {
	return &FooErr{message}
}

func TestHandleError(t *testing.T) {
	var (
		ErrA = errors.New("A")
		ErrB = errors.New("B")
	)

	testCases := []struct {
		name string

		handler grpc.ErrorHandler
		err     error

		expectMessage string
		expectCode    codes.Code
	}{
		{
			name: "DefaultOnly",

			handler: grpc.HandleError(codes.Internal),
			err:     ErrA,

			expectMessage: "A",
			expectCode:    codes.Internal,
		},
		{
			name: "Is",

			handler: grpc.HandleError(codes.Internal).Is(ErrA, codes.InvalidArgument),
			err:     ErrA,

			expectMessage: "A",
			expectCode:    codes.InvalidArgument,
		},
		{
			name: "IsNot",

			handler: grpc.HandleError(codes.Internal).Is(ErrA, codes.InvalidArgument),
			err:     ErrB,

			expectMessage: "B",
			expectCode:    codes.Internal,
		},
		{
			name: "IsW",

			handler: grpc.HandleError(codes.Internal).IsW(ErrA, codes.InvalidArgument, ErrB),
			err:     ErrA,

			expectMessage: "B\nA",
			expectCode:    codes.InvalidArgument,
		},
		{
			name: "IsWNot",

			handler: grpc.HandleError(codes.Internal).IsW(ErrA, codes.InvalidArgument, ErrB),
			err:     ErrB,

			expectMessage: "B",
			expectCode:    codes.Internal,
		},
		{
			name: "IsWF",

			handler: grpc.HandleError(codes.Internal).IsWF(ErrA, codes.InvalidArgument, "Hello %s", "World"),
			err:     ErrA,

			expectMessage: "Hello World",
			expectCode:    codes.InvalidArgument,
		},
		{
			name: "IsWFNot",

			handler: grpc.HandleError(codes.Internal).IsWF(ErrA, codes.InvalidArgument, "Hello %s", "World"),
			err:     ErrB,

			expectMessage: "B",
			expectCode:    codes.Internal,
		},
		{
			name: "As",

			handler: grpc.HandleError(codes.Internal).As(lo.ToPtr(&FooErr{}), codes.InvalidArgument),
			err:     NewFooErr("A"),

			expectMessage: "A",
			expectCode:    codes.InvalidArgument,
		},
		{
			name: "AsNot",

			handler: grpc.HandleError(codes.Internal).As(lo.ToPtr(&FooErr{}), codes.InvalidArgument),
			err:     ErrA,

			expectMessage: "A",
			expectCode:    codes.Internal,
		},
		{
			name: "AsW",

			handler: grpc.HandleError(codes.Internal).AsW(lo.ToPtr(&FooErr{}), codes.InvalidArgument, ErrB),
			err:     NewFooErr("A"),

			expectMessage: "B\nA",
			expectCode:    codes.InvalidArgument,
		},
		{
			name: "AsWNot",

			handler: grpc.HandleError(codes.Internal).AsW(lo.ToPtr(&FooErr{}), codes.InvalidArgument, ErrB),
			err:     ErrA,

			expectMessage: "A",
			expectCode:    codes.Internal,
		},
		{
			name: "AsWF",

			handler: grpc.HandleError(codes.Internal).AsWF(lo.ToPtr(&FooErr{}), codes.InvalidArgument, "Hello %s", "World"),
			err:     NewFooErr("A"),

			expectMessage: "Hello World",
			expectCode:    codes.InvalidArgument,
		},
		{
			name: "AsWFNot",

			handler: grpc.HandleError(codes.Internal).AsWF(lo.ToPtr(&FooErr{}), codes.InvalidArgument, "Hello %s", "World"),
			err:     ErrA,

			expectMessage: "A",
			expectCode:    codes.Internal,
		},
		{
			name: "Test",

			handler: grpc.HandleError(codes.Internal).Test(func(err error) (error, bool) {
				if err.Error() == "A" {
					return status.Errorf(codes.InvalidArgument, err.Error()), true
				}

				return nil, false
			}),
			err: ErrA,

			expectMessage: "A",
			expectCode:    codes.InvalidArgument,
		},
		{
			name: "TestNot",

			handler: grpc.HandleError(codes.Internal).Test(func(err error) (error, bool) {
				if err.Error() == "A" {
					return status.Errorf(codes.InvalidArgument, err.Error()), true
				}

				return nil, false
			}),
			err: ErrB,

			expectMessage: "B",
			expectCode:    codes.Internal,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.handler.Handle(testCase.err)
			require.Error(t, err)

			st, ok := status.FromError(err)
			require.True(t, ok)

			require.Equal(t, testCase.expectMessage, st.Message())
			require.Equal(t, testCase.expectCode, st.Code())
		})
	}
}
