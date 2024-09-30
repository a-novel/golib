package loggers

import (
	"encoding/json"
	"reflect"

	"github.com/rs/zerolog"
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
func (logger *zerologLogger) event(level LogLevel) *zerolog.Event {
	switch level {
	case LogLevelInfo:
		return logger.logger.Info()
	case LogLevelWarning:
		return logger.logger.Warn()
	case LogLevelError:
		return logger.logger.Error()
	case LogLevelFatal:
		return logger.logger.Fatal()
	default:
		return logger.logger.Info()
	}
}

// AsZeroLoggable checks if an interface implements the ZerologMessage interface. It returns the underlying
// implementation on success.
func AsZeroLoggable(data interface{}) (ZerologMessage, bool) {
	zerologMessage, ok := data.(ZerologMessage)
	return zerologMessage, ok
}

func (logger *zerologLogger) printStruct(event *zerolog.Event, data interface{}) {
	mrsh, err := json.Marshal(data)
	// This should not happen but whatever. In any case it would be an error on our side.
	if err != nil {
		event.Interface("data", data).Msg("")
		return
	}

	var out map[string]interface{}
	if err = json.Unmarshal(mrsh, &out); err != nil {
		event.Interface("data", data).Msg("")
		return
	}

	for key, value := range out {
		event = event.Interface(key, value)
	}

	event.Msg("")
}

func (logger *zerologLogger) printMap(event *zerolog.Event, data interface{}) {
	keyType := reflect.TypeOf(data).Key().Kind()

	if keyType != reflect.String {
		event.Interface("data", data).Msg("")
		return
	}

	// Get the reflect value to loop through the data actual content.
	v := reflect.ValueOf(data)

	for _, key := range v.MapKeys() {
		// Push the value to the zerolog event.
		event = event.Interface(key.String(), v.MapIndex(key).Interface())
	}

	event.Msg("")
}

func (logger *zerologLogger) Log(level LogLevel, data interface{}) {
	// Don't log empty messages, do a no-op instead.
	if data == nil {
		return
	}

	// Generate a new empty event from log level.
	event := logger.event(level)

	// If the log data implements its own ZeroLog method, use it.
	// This is only available for top-level messages.
	customMarshaller, ok := AsZeroLoggable(data)
	if ok {
		customMarshaller.ZeroLog(event)
		return
	}

	// Use reflect to make the best use of zerolog capabilities, depending on the message type.
	t := reflect.TypeOf(data)
	switch t.Kind() {
	case reflect.String:
		// Strings are directly sent as messages.
		event.Msg(data.(string))
	case
		reflect.Int,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		// Numbers are stringified, then sent as messages.
		event.Msgf("%v", data)
	case reflect.Struct:
		logger.printStruct(event, data)
	case reflect.Map:
		logger.printMap(event, data)
	default:
		// No effort at this point, just add the data as an interface then send it.
		event.Interface("data", data).Msg("")
	}
}

// NewZeroLog creates a new JSON that writes to a zerolog.Logger.
//
// It can make use of messages that implement the ZerologMessage interface.
func NewZeroLog(logger zerolog.Logger) JSON {
	return &zerologLogger{logger: logger}
}
