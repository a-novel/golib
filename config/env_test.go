package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/config"
)

func TestEnvStringParser(t *testing.T) {
	t.Setenv("foo", "foo:value")

	require.Equal(t, "foo:value", config.LoadEnv(os.Getenv("foo"), "foo:default", config.StringParser))
	require.Equal(t, "bar:default", config.LoadEnv(os.Getenv("bar"), "bar:default", config.StringParser))
}

func TestEnvInt64Parser(t *testing.T) {
	t.Setenv("foo", "123")
	t.Setenv("bar", "bar:value")

	require.Equal(t, int64(123), config.LoadEnv(os.Getenv("foo"), 321, config.Int64Parser))
	require.Equal(t, int64(321), config.LoadEnv(os.Getenv("bar"), 321, config.Int64Parser))
	require.Equal(t, int64(321), config.LoadEnv(os.Getenv("qux"), 321, config.Int64Parser))
}

func TestEnvInt32Parser(t *testing.T) {
	t.Setenv("foo", "123")
	t.Setenv("bar", "bar:value")

	require.Equal(t, int32(123), config.LoadEnv(os.Getenv("foo"), 321, config.Int32Parser))
	require.Equal(t, int32(321), config.LoadEnv(os.Getenv("bar"), 321, config.Int32Parser))
	require.Equal(t, int32(321), config.LoadEnv(os.Getenv("qux"), 321, config.Int32Parser))
}

func TestEnvInt16Parser(t *testing.T) {
	t.Setenv("foo", "123")
	t.Setenv("bar", "bar:value")

	require.Equal(t, int16(123), config.LoadEnv(os.Getenv("foo"), 321, config.Int16Parser))
	require.Equal(t, int16(321), config.LoadEnv(os.Getenv("bar"), 321, config.Int16Parser))
	require.Equal(t, int16(321), config.LoadEnv(os.Getenv("qux"), 321, config.Int16Parser))
}

func TestEnvInt8Parser(t *testing.T) {
	t.Setenv("foo", "123")
	t.Setenv("bar", "bar:value")

	require.Equal(t, int8(123), config.LoadEnv(os.Getenv("foo"), 64, config.Int8Parser))
	require.Equal(t, int8(64), config.LoadEnv(os.Getenv("bar"), 64, config.Int8Parser))
	require.Equal(t, int8(64), config.LoadEnv(os.Getenv("qux"), 64, config.Int8Parser))
}

func TestEnvIntParser(t *testing.T) {
	t.Setenv("foo", "123")
	t.Setenv("bar", "bar:value")

	require.Equal(t, 123, config.LoadEnv(os.Getenv("foo"), 321, config.IntParser))
	require.Equal(t, 321, config.LoadEnv(os.Getenv("bar"), 321, config.IntParser))
	require.Equal(t, 321, config.LoadEnv(os.Getenv("qux"), 321, config.IntParser))
}

func TestEnvUint64Parser(t *testing.T) {
	t.Setenv("foo", "123")
	t.Setenv("bar", "bar:value")

	require.Equal(t, uint64(123), config.LoadEnv(os.Getenv("foo"), 321, config.Uint64Parser))
	require.Equal(t, uint64(321), config.LoadEnv(os.Getenv("bar"), 321, config.Uint64Parser))
	require.Equal(t, uint64(321), config.LoadEnv(os.Getenv("qux"), 321, config.Uint64Parser))
}

func TestEnvUint32Parser(t *testing.T) {
	t.Setenv("foo", "123")
	t.Setenv("bar", "bar:value")

	require.Equal(t, uint32(123), config.LoadEnv(os.Getenv("foo"), 321, config.Uint32Parser))
	require.Equal(t, uint32(321), config.LoadEnv(os.Getenv("bar"), 321, config.Uint32Parser))
	require.Equal(t, uint32(321), config.LoadEnv(os.Getenv("qux"), 321, config.Uint32Parser))
}

func TestEnvUint16Parser(t *testing.T) {
	t.Setenv("foo", "123")
	t.Setenv("bar", "bar:value")

	require.Equal(t, uint16(123), config.LoadEnv(os.Getenv("foo"), 321, config.Uint16Parser))
	require.Equal(t, uint16(321), config.LoadEnv(os.Getenv("bar"), 321, config.Uint16Parser))
	require.Equal(t, uint16(321), config.LoadEnv(os.Getenv("qux"), 321, config.Uint16Parser))
}

func TestEnvUint8Parser(t *testing.T) {
	t.Setenv("foo", "123")
	t.Setenv("bar", "bar:value")

	require.Equal(t, uint8(123), config.LoadEnv(os.Getenv("foo"), 64, config.Uint8Parser))
	require.Equal(t, uint8(64), config.LoadEnv(os.Getenv("bar"), 64, config.Uint8Parser))
	require.Equal(t, uint8(64), config.LoadEnv(os.Getenv("qux"), 64, config.Uint8Parser))
}

func TestEnvUintParser(t *testing.T) {
	t.Setenv("foo", "123")
	t.Setenv("bar", "bar:value")

	require.Equal(t, uint(123), config.LoadEnv(os.Getenv("foo"), 321, config.UintParser))
	require.Equal(t, uint(321), config.LoadEnv(os.Getenv("bar"), 321, config.UintParser))
	require.Equal(t, uint(321), config.LoadEnv(os.Getenv("qux"), 321, config.UintParser))
}

func TestEnvBoolParser(t *testing.T) {
	t.Setenv("foo", "false")
	t.Setenv("bar", "bar:value")

	require.False(t, config.LoadEnv(os.Getenv("foo"), true, config.BoolParser))
	require.True(t, config.LoadEnv(os.Getenv("bar"), true, config.BoolParser))
	require.True(t, config.LoadEnv(os.Getenv("qux"), true, config.BoolParser))
}

func TestEnvFloat64Parser(t *testing.T) {
	t.Setenv("foo", "123")
	t.Setenv("bar", "bar:value")

	require.InEpsilon(t, float64(123), config.LoadEnv(os.Getenv("foo"), 321, config.Float64Parser), 0.0001)
	require.InEpsilon(t, float64(321), config.LoadEnv(os.Getenv("bar"), 321, config.Float64Parser), 0.0001)
	require.InEpsilon(t, float64(321), config.LoadEnv(os.Getenv("qux"), 321, config.Float64Parser), 0.0001)
}

func TestEnvFloat32Parser(t *testing.T) {
	t.Setenv("foo", "123")
	t.Setenv("bar", "bar:value")

	require.InEpsilon(t, float32(123), config.LoadEnv(os.Getenv("foo"), 321, config.Float32Parser), 0.0001)
	require.InEpsilon(t, float32(321), config.LoadEnv(os.Getenv("bar"), 321, config.Float32Parser), 0.0001)
	require.InEpsilon(t, float32(321), config.LoadEnv(os.Getenv("qux"), 321, config.Float32Parser), 0.0001)
}

func TestEnvDurationParser(t *testing.T) {
	t.Setenv("foo", "5s")
	t.Setenv("bar", "bar:value")

	require.Equal(t, 5*time.Second, config.LoadEnv(os.Getenv("foo"), time.Second, config.DurationParser))
	require.Equal(t, time.Second, config.LoadEnv(os.Getenv("bar"), time.Second, config.DurationParser))
	require.Equal(t, time.Second, config.LoadEnv(os.Getenv("qux"), time.Second, config.DurationParser))
}

func TestEnvTimeParser(t *testing.T) {
	t.Setenv("foo", "2020-01-02T15:04:05Z")
	t.Setenv("bar", "bar:value")

	now := time.Now()
	custom := time.Date(2020, 1, 2, 15, 4, 5, 0, time.UTC)

	require.Equal(t, custom, config.LoadEnv(os.Getenv("foo"), now, config.TimeParser))
	require.Equal(t, now, config.LoadEnv(os.Getenv("bar"), now, config.TimeParser))
	require.Equal(t, now, config.LoadEnv(os.Getenv("qux"), now, config.TimeParser))
}

func TestEnvJSONMapParser(t *testing.T) {
	t.Setenv("foo", `{"key":"super-value"}`)
	t.Setenv("bar", "bar:value")

	basic := map[string]any{"key": "value"}
	custom := map[string]any{"key": "super-value"}

	require.Equal(t, custom, config.LoadEnv(os.Getenv("foo"), basic, config.JSONMapParser))
	require.Equal(t, basic, config.LoadEnv(os.Getenv("bar"), basic, config.JSONMapParser))
	require.Equal(t, basic, config.LoadEnv(os.Getenv("qux"), basic, config.JSONMapParser))
}

func TestEnvJSONSliceParser(t *testing.T) {
	t.Setenv("foo", `[{"key":"super-value"}]`)
	t.Setenv("bar", "bar:value")

	basic := []any{map[string]any{"key": "value"}}
	custom := []any{map[string]any{"key": "super-value"}}

	require.Equal(t, custom, config.LoadEnv(os.Getenv("foo"), basic, config.JSONSliceParser))
	require.Equal(t, basic, config.LoadEnv(os.Getenv("bar"), basic, config.JSONSliceParser))
	require.Equal(t, basic, config.LoadEnv(os.Getenv("qux"), basic, config.JSONSliceParser))
}

func TestEnvStringSliceParser(t *testing.T) {
	t.Setenv("foo", "foo,bar,baz")

	basic := []string{"foo", "bar", "qux"}
	custom := []string{"foo", "bar", "baz"}

	require.Equal(t, custom, config.LoadEnv(os.Getenv("foo"), basic, config.SliceParser(config.StringParser)))
	require.Equal(t, basic, config.LoadEnv(os.Getenv("bar"), basic, config.SliceParser(config.StringParser)))
}

func TestEnvEnumParser(t *testing.T) {
	t.Setenv("foo", "bar")

	require.Equal(
		t,
		"bar",
		config.LoadEnv(os.Getenv("foo"), "invalid", config.EnumParser(config.StringParser, "foo", "bar")),
	)
	require.Equal(
		t,
		"invalid",
		config.LoadEnv(os.Getenv("foo"), "invalid", config.EnumParser(config.StringParser, "foo", "baz")),
	)
}
