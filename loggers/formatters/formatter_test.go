package formatters_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/a-novel/golib/loggers"
	"github.com/a-novel/golib/loggers/formatters"
	formattersmocks "github.com/a-novel/golib/loggers/formatters/mocks"
	loggersmocks "github.com/a-novel/golib/loggers/mocks"
)

func TestConsoleFormatter(t *testing.T) {
	t.Run("Render", func(t *testing.T) {
		testCases := []struct {
			name string

			content string
			level   loggers.LogLevel
		}{
			{
				name:    "Info",
				content: "foo",
				level:   loggers.LogLevelInfo,
			},
			{
				name:    "Warning",
				content: "foo",
				level:   loggers.LogLevelWarning,
			},
			{
				name:    "Error",
				content: "foo",
				level:   loggers.LogLevelError,
			},
			{
				name:    "Fatal",
				content: "foo",
				level:   loggers.LogLevelFatal,
			},
		}

		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				content := formattersmocks.NewMockLogContent(t)
				logger := loggersmocks.NewMockConsole(t)

				content.On("RenderConsole").Return(tt.content)
				logger.On("Log", tt.level, tt.content)

				formatter := formatters.NewConsoleFormatter(logger, false)

				formatter.Log(content, tt.level)

				content.AssertExpectations(t)
				logger.AssertExpectations(t)
			})
		}
	})

	t.Run("RenderDynamic", func(t *testing.T) {
		testCases := []struct {
			name string

			content string
			level   loggers.LogLevel
		}{
			{
				name:    "Info",
				content: "foo",
				level:   loggers.LogLevelInfo,
			},
			{
				name:    "Warning",
				content: "foo",
				level:   loggers.LogLevelWarning,
			},
			{
				name:    "Error",
				content: "foo",
				level:   loggers.LogLevelError,
			},
			{
				name:    "Fatal",
				content: "foo",
				level:   loggers.LogLevelFatal,
			},
		}

		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				content := formattersmocks.NewMockLogDynamicContent(t)
				logger := loggersmocks.NewMockConsole(t)

				content.
					On(
						"RenderConsoleDynamic",
						mock.AnythingOfType(reflect.FuncOf(
							[]reflect.Type{reflect.TypeOf("")},
							nil,
							false,
						).String()),
					).
					Run(func(args mock.Arguments) {
						fn := args[0].(func(string))
						fn(tt.content)
					}).
					Return(content)

				logger.On("Log", tt.level, tt.content)

				formatter := formatters.NewConsoleFormatter(logger, false)

				formatter.Log(content, tt.level)

				content.AssertExpectations(t)
				logger.AssertExpectations(t)
			})
		}
	})

	t.Run("RenderDynamicUnderStaticEnvironment", func(t *testing.T) {
		content := formattersmocks.NewMockLogDynamicContent(t)
		logger := loggersmocks.NewMockConsole(t)

		content.On("RenderConsole").Return("foo")
		logger.On("Log", loggers.LogLevelInfo, "foo")

		formatter := formatters.NewConsoleFormatter(logger, true)

		formatter.Log(content, loggers.LogLevelInfo)

		content.AssertExpectations(t)
		logger.AssertExpectations(t)
	})
}

func TestNewJSONFormatter(t *testing.T) {
	t.Run("Render", func(t *testing.T) {
		testCases := []struct {
			name string

			content interface{}
			level   loggers.LogLevel
		}{
			{
				name:    "Info",
				content: map[string]interface{}{"foo": "bar"},
				level:   loggers.LogLevelInfo,
			},
			{
				name:    "Warning",
				content: map[string]interface{}{"foo": "bar"},
				level:   loggers.LogLevelWarning,
			},
			{
				name:    "Error",
				content: map[string]interface{}{"foo": "bar"},
				level:   loggers.LogLevelError,
			},
			{
				name:    "Fatal",
				content: map[string]interface{}{"foo": "bar"},
				level:   loggers.LogLevelFatal,
			},
		}

		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				content := formattersmocks.NewMockLogContent(t)
				logger := loggersmocks.NewMockJSON(t)

				content.On("RenderJSON").Return(tt.content)
				logger.On("Log", tt.level, tt.content)

				formatter := formatters.NewJSONFormatter(logger)

				formatter.Log(content, tt.level)

				content.AssertExpectations(t)
				logger.AssertExpectations(t)
			})
		}
	})

	t.Run("IgnoreDynamicContent", func(t *testing.T) {
		content := formattersmocks.NewMockLogDynamicContent(t)
		logger := loggersmocks.NewMockJSON(t)

		content.On("RenderJSON").Return(map[string]interface{}{"foo": "bar"})
		logger.On("Log", loggers.LogLevelInfo, map[string]interface{}{"foo": "bar"})

		formatter := formatters.NewJSONFormatter(logger)

		formatter.Log(content, loggers.LogLevelInfo)

		content.AssertExpectations(t)
		logger.AssertExpectations(t)
	})
}
