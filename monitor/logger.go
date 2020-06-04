package monitor


import (
	"context"
	"net/http"

	gcontext "github.com/cultureamp/glamplify/context"
	"github.com/cultureamp/glamplify/log"
)

// Logger
type Logger struct {
	internal *log.Logger
}


// New creates a *Logger with optional fields. Useful for when you want to add a field to all subsequent logging calls eg. request_id, etc.
func New(rsFields gcontext.RequestScopedFields, fields ...log.Fields) *Logger {
	return newLogger(rsFields, fields...)
}

// NewFromCtx creates a new logger from a context, which should contain RequestScopedFields.
// If the context does not contain then, then this method will NOT add them in.
func NewFromCtx(ctx context.Context, fields ...log.Fields) *Logger {
	rsFields, _ := gcontext.GetRequestScopedFields(ctx)
	return New(rsFields, fields...)
}

// NewFromRequest creates a new logger from a http.Request, which should contain RequestScopedFields.
// If the context does not contain then, then this method will NOT add them in.
func NewFromRequest(r *http.Request, fields ...log.Fields) *Logger {
	return NewFromCtx(r.Context(), fields...)
}

func newLogger(rsFields gcontext.RequestScopedFields, fields ...log.Fields) *Logger {

	writer := NewWriter()
	internal := log.NewWitCustomWriter(rsFields, writer, fields...)

	logger := &Logger{
		internal: internal,
	}
	return logger
}

// Debug writes a write message with optional types to the underlying standard writer.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// Use snake_case keys and lower case values if possible.
func Debug(rsFields gcontext.RequestScopedFields, event string, fields ...log.Fields) {
	defaultLogger := New(rsFields)
	defaultLogger.Debug(event, fields...)
}

// Debug writes a write message with optional types to the underlying standard writer.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Debug(event string, fields ...log.Fields) {
	logger.internal.Debug(event, fields...)
}

// Info writes a message with optional types to the underlying standard writer.
// Useful for normal tracing that should be captured during standard operating behaviour.
// Use snake_case keys and lower case values if possible.
func Info(rsFields gcontext.RequestScopedFields, event string, fields ...log.Fields) {
	defaultLogger := New(rsFields)
	defaultLogger.Info(event, fields...)
}

// Info writes a message with optional types to the underlying standard writer.
// Useful for normal tracing that should be captured during standard operating behaviour.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Info(event string, fields ...log.Fields) {
	logger.internal.Info(event, fields...)
}

// Warn writes a message with optional types to the underlying standard writer.
// Useful for unusual but recoverable tracing that should be captured during standard operating behaviour.
// Use snake_case keys and lower case values if possible.
func Warn(rsFields gcontext.RequestScopedFields, event string, fields ...log.Fields) {
	defaultLogger := New(rsFields)
	defaultLogger.Warn(event, fields...)
}

// Warn writes a message with optional types to the underlying standard writer.
// Useful for unusual but recoverable tracing that should be captured during standard operating behaviour.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Warn(event string, fields ...log.Fields) {
	logger.internal.Warn(event, fields...)
}

// Error writes a error message with optional types to the underlying standard writer.
// Useful to trace errors that are usually not recoverable. These should always be logged.
// Use snake_case keys and lower case values if possible.
func Error(rsFields gcontext.RequestScopedFields, event string, err error, fields ...log.Fields) {
	defaultLogger := New(rsFields)
	defaultLogger.Error(event, err, fields...)
}

// Error writes a error message with optional types to the underlying standard writer.
// Useful to trace errors that are usually not recoverable. These should always be logged.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Error(event string, err error, fields ...log.Fields) {
	logger.internal.Error(event, err, fields...)
}

// Fatal writes a error message with optional types to the underlying standard writer and then calls panic!
// Panic will terminate the current go routine.
// Useful to trace catastrophic errors that are not recoverable. These should always be logged.
// Use snake_case keys and lower case values if possible.
func Fatal(rsFields gcontext.RequestScopedFields, event string, err error, fields ...log.Fields) {
	defaultLogger := New(rsFields)
	defaultLogger.Fatal(event, err, fields...)
}

// Fatal writes a error message with optional types to the underlying standard writer and then calls panic!
// Panic will terminate the current go routine.
// Useful to trace catastrophic errors that are not recoverable. These should always be logged.
// Use snake_case keys and lower case values if possible.
func (logger Logger) Fatal(event string, err error, fields ...log.Fields) {
	logger.internal.Fatal(event, err, fields...)
}

// Event method uses expressive syntax format: logger.Event("event_name").Fields(fields...).Info("message")
func (logger Logger) Event(event string) *Segment {

	return &Segment{
		logger: logger,
		event: event,
		fields: log.Fields{},
	}
}

