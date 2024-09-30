package loggers_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/golib/loggers"
	"github.com/a-novel/golib/testutils"
)

type customZeroLogStruct struct{}

func (c customZeroLogStruct) ZeroLog(e *zerolog.Event) {
	e.Str("foo", "bar")
	e.Msg("Hell Yeah")
}

func TestZeroLog(t *testing.T) {
	testData := []struct {
		name string

		level loggers.LogLevel
		msg   interface{}

		expect interface{}
	}{
		{
			name:  "Level/Info",
			level: loggers.LogLevelInfo,
			msg:   "foo",

			expect: map[string]interface{}{
				"level":   "info",
				"message": "foo",
			},
		},
		{
			name:  "Level/Warn",
			level: loggers.LogLevelWarning,
			msg:   "foo",

			expect: map[string]interface{}{
				"level":   "warn",
				"message": "foo",
			},
		},
		{
			name:  "Level/Error",
			level: loggers.LogLevelError,
			msg:   "foo",

			expect: map[string]interface{}{
				"level":   "error",
				"message": "foo",
			},
		},
		{
			name:  "Level/Unknown",
			level: loggers.LogLevel("unknown"),
			msg:   "foo",

			expect: map[string]interface{}{
				"level":   "info",
				"message": "foo",
			},
		},

		{
			name:  "DataType/Nil",
			level: loggers.LogLevelInfo,
		},
		{
			name:  "DataType/String",
			level: loggers.LogLevelInfo,
			msg:   "foo",

			expect: map[string]interface{}{
				"level":   "info",
				"message": "foo",
			},
		},
		{
			name:  "DataType/Int",
			level: loggers.LogLevelInfo,
			msg:   10,

			expect: map[string]interface{}{
				"level":   "info",
				"message": "10",
			},
		},
		{
			name:  "DataType/StringMap",
			level: loggers.LogLevelInfo,
			msg:   map[string]interface{}{"foo": "bar", "hello": "kitty"},

			expect: map[string]interface{}{
				"level": "info",
				"foo":   "bar",
				"hello": "kitty",
			},
		},
		{
			name:  "DataType/StringMap/MessageKey",
			level: loggers.LogLevelInfo,
			msg:   map[string]interface{}{"message": "foo bar", "hello": "kitty"},

			expect: map[string]interface{}{
				"level":   "info",
				"message": "foo bar",
				"hello":   "kitty",
			},
		},
		{
			name:  "DataType/Map",
			level: loggers.LogLevelInfo,
			msg:   map[int]interface{}{1: "item1", 2: "item2"},

			expect: map[string]interface{}{
				"level": "info",
				"data":  map[string]interface{}{"1": "item1", "2": "item2"},
			},
		},
		{
			name:  "DataType/Struct",
			level: loggers.LogLevelInfo,
			msg: struct {
				Foo   string `json:"foo"`
				Hello string `json:"hello"`
			}{
				Foo:   "bar",
				Hello: "kitty",
			},

			expect: map[string]interface{}{
				"level": "info",
				"foo":   "bar",
				"hello": "kitty",
			},
		},
		{
			name:  "DataType/Other",
			level: loggers.LogLevelInfo,
			msg:   []string{"item1", "item2"},

			expect: map[string]interface{}{
				"level": "info",
				"data":  []interface{}{"item1", "item2"},
			},
		},
		{
			name:  "DataType/CustomZeroLoggable",
			level: loggers.LogLevelInfo,
			msg:   new(customZeroLogStruct),

			expect: map[string]interface{}{
				"level":   "info",
				"foo":     "bar",
				"message": "Hell Yeah",
			},
		},
	}

	for _, tc := range testData {
		t.Run(tc.name, func(t *testing.T) {
			w, capture, err := testutils.CreateSTDCapture(t)
			require.NoError(t, err)

			srcLogger := zerolog.New(w)
			logger := loggers.NewZeroLog(srcLogger)

			logger.Log(tc.level, tc.msg)

			// We unmarshal the log, rather than comparing strings, because the output is not deterministic.
			rawCapture := capture()
			if tc.msg == nil {
				require.Empty(t, rawCapture)
				return
			}

			require.NotEmpty(t, rawCapture)

			var captured interface{}
			require.NoError(t, json.Unmarshal([]byte(rawCapture), &captured))

			require.Equal(t, tc.expect, captured)
		})
	}
}

func TestZeroLogPanic(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			srcLogger := zerolog.New(os.Stdout)
			logger := loggers.NewZeroLog(srcLogger)
			logger.Log(loggers.LogLevelFatal, "oopsie doopsie")
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.False(t, res.Success)
			require.Equal(t, res.STDOut, "{\"level\":\"fatal\",\"message\":\"oopsie doopsie\"}\n")
		},
	})
}
