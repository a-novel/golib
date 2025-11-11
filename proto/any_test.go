package proto_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/proto"
)

func TestProtoAnyConversion(t *testing.T) {
	t.Parallel()

	anyValue := map[string]any{
		"message": "hello world",
	}

	toProto, err := proto.InterfaceToProtoAny(anyValue)
	require.NoError(t, err)

	fromProto, err := proto.ProtoAnyToInterface(toProto)
	require.NoError(t, err)

	require.Equal(t, anyValue, fromProto)
}
