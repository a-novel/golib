package logger

import (
	"github.com/rs/zerolog"
	"reflect"
)

// ZerologMessage is a custom interface messages can use to customize the zerolog event.
//
// If a message implements this interface, the logger will rely on the ZeroLog method to produce a message, rather than
// doing it itself.
//
// The ZeroLog method is responsible for triggering the message with e.Msg or e.Msgf.
//
// This method is only called for top-level messages. Nested messages will use the default serialization behavior.
type ZerologMessage interface {
	// ZeroLog shapes, and sends a customized message using zerolog event.
	ZeroLog(e *zerolog.Event)
}

// Implements JSON using the zerolog.Logger interface.
type zerologLogger struct {
	logger zerolog.Logger
}

// Create a custom zerolog.Event depending on the log level.
func (l *zerologLogger) event(level LogLevel) *zerolog.Event {
	switch level {
	case LogLevelInfo:
		return l.logger.Info()
	case LogLevelWarning:
		return l.logger.Warn()
	case LogLevelError:
		return l.logger.Error()
	case LogLevelFatal:
		return l.logger.Fatal()
	default:
		return l.logger.Info()
	}
}

// AsZeroLoggable checks if an interface implements the ZerologMessage interface. It returns the underlying
// implementation on success.
func AsZeroLoggable(data interface{}) (ZerologMessage, bool) {
	zm, ok := data.(ZerologMessage)
	return zm, ok
}

func (l *zerologLogger) Log(level LogLevel, data interface{}) {
	// Don't log empty messages, do a no-op instead.
	if data == nil {
		return
	}

	// Generate a new empty event from log level.
	e := l.event(level)

	// If the log data implements its own ZeroLog method, use it.
	// This is only available for top-level messages.
	customMarshaller, ok := AsZeroLoggable(data)
	if ok {
		customMarshaller.ZeroLog(e)
		return
	}

	// Use reflect to make the best use of zerolog capabilities, depending on the message type.
	t := reflect.TypeOf(data)
	switch t.Kind() {
	case reflect.String:
		// Strings are directly sent as messages.
		e.Msg(data.(string))
	case
		reflect.Int,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		// Numbers are stringified, then sent as messages.
		e.Msgf("%v", data)
	case reflect.Map:
		// Only string maps are natively supported. Other maps will fall through the default behavior.
		if t.Key().Kind() == reflect.String {
			// Get the reflect value to loop through the data actual content.
			v := reflect.ValueOf(data)
			// Catch the message key to give it a special treatment (it will be logged as the message, rather than
			// a standard key).
			msg := ""

			for _, key := range v.MapKeys() {
				// Save message key in a separate variable, and skip it from the loop.
				if key.String() == "message" && v.MapIndex(key).Kind() == reflect.String {
					msg = v.MapIndex(key).String()
					continue
				}

				// Push the value to the zerolog event.
				e = e.Interface(key.String(), v.MapIndex(key).Interface())
			}

			e.Msg(msg)
			return
		}

		// Unhandled maps to default.
		fallthrough
	default:
		// No effort at this point, just add the data as an interface then send it.
		e.Interface("data", data).Msg("")
	}
}

// NewZerolog creates a new JSON that writes to a zerolog.Logger.
//
// It can make use of messages that implement the ZerologMessage interface.
func NewZerolog(logger zerolog.Logger) JSON {
	return &zerologLogger{logger: logger}
}
