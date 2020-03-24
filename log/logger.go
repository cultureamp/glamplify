package log

import (
	"context"
	"github.com/cultureamp/glamplify/constants"
	"strings"
)

// Logger allows you to set types that can be re-used for subsequent log event. Useful for setting username, requestid etc for a Http Web Request.
type Logger struct {
	ctx       context.Context
	writer    *FieldWriter
	fields    Fields
	defValues *DefaultValues
}

var (
	internal = newWriter(func(conf *config) {
	})
)

// FromScope lets you add types to a scoped writer. Useful for Http Web Request where you want to track user, requestid, etc.
func FromScope(ctx context.Context, fields ...Fields) *Logger {
	return newLogger(ctx, internal, fields...)
}

func newLogger(ctx context.Context, writer *FieldWriter, fields ...Fields) *Logger {
	merged := Fields{}
	merged = merged.Merge(fields...)
	scope := &Logger{
		ctx:    ctx,
		writer: writer,
		fields: merged,
	}
	scope.defValues = newDefaultValues()
	return scope
}

// Debug writes a debug message with optional types to the underlying standard writer.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// Use lower-case keys and values if possible.
func (logger Logger) Debug(event string, fields ...Fields) {
	meta := logger.defValues.getDefaults(logger.ctx, event, constants.DebugSevLogValue)
	merged := logger.fields.Merge(fields...)
	logger.writer.debug(event, meta, merged)
}

// Info writes a message with optional types to the underlying standard writer.
// Useful for normal tracing that should be captured during standard operating behaviour.
// Use lower-case keys and values if possible.
func (logger Logger) Info(event string, fields ...Fields) {
	meta := logger.defValues.getDefaults(logger.ctx, event, constants.InfoSevLogValue)
	merged := logger.fields.Merge(fields...)
	logger.writer.info(event, meta, merged)
}

// Warn writes a message with optional types to the underlying standard writer.
// Useful for unusual but recoverable tracing that should be captured during standard operating behaviour.
// Use lower-case keys and values if possible.
func (logger Logger) Warn(event string, fields ...Fields) {
	meta := logger.defValues.getDefaults(logger.ctx, event, constants.WarnSevLogValue)
	merged := logger.fields.Merge(fields...)
	logger.writer.warn(event, meta, merged)
}

// Error writes a error message with optional types to the underlying standard writer.
// Useful to trace errors that are usually not recoverable. These should always be logged.
// Use lower-case keys and values if possible.
func (logger Logger) Error(err error, fields ...Fields) {
	event := strings.TrimSpace(err.Error())
	meta := logger.defValues.getDefaults(logger.ctx, event, constants.ErrorSevLogValue)
	meta = logger.defValues.getErrorDefaults(err, meta)
	merged := logger.fields.Merge(fields...)
	logger.writer.error(event, meta, merged)
}

// Fatal writes a error message with optional types to the underlying standard writer and then calls panic!
// Panic will terminate the current go routine.
// Useful to trace catastrophic errors that are not recoverable. These should always be logged.
// Use lower-case keys and values if possible.
func (logger Logger) Fatal(err error, fields ...Fields) {
	event := strings.TrimSpace(err.Error())
	meta := logger.defValues.getDefaults(logger.ctx, event, constants.FatalSevLogValue)
	meta = logger.defValues.getErrorDefaults(err, meta)
	merged := logger.fields.Merge(fields...)
	logger.writer.fatal(event, meta, merged)
}
