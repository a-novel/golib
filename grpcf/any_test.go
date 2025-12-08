package grpcf_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/golib/grpcf"
)

func TestProtoAnyConversion(t *testing.T) {
	t.Parallel()

	anyValue := map[string]any{
		"message": "hello world",
	}

	toProto, err := grpcf.InterfaceToProtoAny(anyValue)
	require.NoError(t, err)

	fromProto, err := grpcf.ProtoAnyToInterface(toProto)
	require.NoError(t, err)

	require.Equal(t, anyValue, fromProto)
}
