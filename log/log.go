package log

import (
	"io"
	"os"
	"strings"
	"sync"
)

// RFC3339Milli is the standard RFC3339 format with added milliseconds
const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

// Config for setting initial values for Logger
type Config struct {
	Output            io.Writer
	TimeFormat        string
	debugForwardLogTo string
	printForwardLogTo string
	errorForwardLogTo string
}

// FieldLogger wraps the standard library logger and add structured field as quoted key value pairs
type FieldLogger struct {
	mutex             sync.Mutex
	output            io.Writer
	timeFormat        string
	debugForwardLogTo string
	printForwardLogTo string
	errorForwardLogTo string
}

// So that you don't even need to create a new logger
var (
		internal = New(func(conf *Config) {
		conf.debugForwardLogTo = "none"
		conf.printForwardLogTo = "splunk"
		conf.errorForwardLogTo = "splunk"
	})
)

// New creates a new FieldLogger. The optional configure func lets you set values on the underlying standard logger.
// eg. SetOutput
func New(configure ...func(*Config)) *FieldLogger { // https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

	logger := &FieldLogger{}
	conf := Config{
		Output:            os.Stdout,
		TimeFormat:        RFC3339Milli,
		debugForwardLogTo: "none",
		printForwardLogTo: "splunk",
		errorForwardLogTo: "splunk",
	}
	for _, config := range configure {
		config(&conf)
	}

	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	logger.output = conf.Output
	logger.timeFormat = conf.TimeFormat
	logger.debugForwardLogTo = conf.debugForwardLogTo
	logger.printForwardLogTo = conf.printForwardLogTo
	logger.errorForwardLogTo = conf.errorForwardLogTo

	return logger
}

// WithScope lets you add field to a scoped logger. Useful for Http Web Request where you want to track user, requestid, etc.
func WithScope(fields Fields) *Scope {
	return newScope(internal, fields)
}

// WithScope lets you add field to a scoped logger. Useful for Http Web Request where you want to track user, requestid, etc.
func (logger *FieldLogger) WithScope(fields Fields) *Scope {
	return newScope(logger, fields)
}

// Debug writes a debug message with optional field to the underlying standard logger.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// All field values will be automatically quoted (keys will not be).
// Debug adds field {level="debug", time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "field message"
// Use lower-case keys and values if possible.
func Debug(message string, fields ...Fields) {
	internal.Debug(message, fields...)
}

// Debug writes a debug message with optional field to the underlying standard logger.
// Useful for adding detailed tracing that you don't normally want to appear, but turned on
// when hunting down incorrect behaviour.
// All field values will be automatically quoted (keys will not be).
// Debug adds field {level="debug", time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "field message"
// Use lower-case keys and values if possible.
func (logger *FieldLogger) Debug(message string, fields ...Fields) {
	meta := Fields{
		HOST:     hostName(),
		MESSAGE:  message,
		PID:      processID(),
		PROCESS:  processName(),
		SEVERITY: DEBUG_SEV,
		TIME:     timeNow(logger.timeFormat),
		FORWARD:  logger.debugForwardLogTo,
	}

	merged := meta.merge(fields...)
	str := merged.serialize()
	logger.write(str)
}

// Print writes a message with optional field to the underlying standard logger.
// Useful to normal tracing that should be captured during standard operating behaviour.
// All field values will be automatically quoted (keys will not be).
// Debug adds field {time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "field message"
// Use lower-case keys and values if possible.
func Print(message string, fields ...Fields) {
	internal.Print(message, fields...)
}

// Print writes a message with optional field to the underlying standard logger.
// Useful to normal tracing that should be captured during standard operating behaviour.
// All field values will be automatically quoted (keys will not be).
// Debug adds field {time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "field message"
// Use lower-case keys and values if possible.
func (logger *FieldLogger) Print(message string, fields ...Fields) {
	meta := Fields{
		MESSAGE:  message,
		SEVERITY: INFO_SEV,
		TIME:     timeNow(logger.timeFormat),
		FORWARD:  logger.printForwardLogTo,
	}

	merged := meta.merge(fields...)
	str := merged.serialize()
	logger.write(str)
}

// Error writes a error message with optional field to the underlying standard logger.
// Useful to trace errors should be captured always.
// All field values will be automatically quoted (keys will not be).
// Debug adds field {level="error", time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "field message"
// Use lower-case keys and values if possible.
func Error(err error, fields ...Fields) {
	internal.Error(err, fields...)
}

// Error writes a error message with optional field to the underlying standard logger.
// Useful to trace errors should be captured always.
// All field values will be automatically quoted (keys will not be).
// Debug adds field {level="error", time="2006-01-02T15:04:05Z07:00"}
// and prints output in the format "field message"
// Use lower-case keys and values if possible.
func (logger *FieldLogger) Error(err error, fields ...Fields) {
	meta := Fields{
		ARCHITECTURE: targetArch(),
		ERROR:        strings.TrimSpace(err.Error()),
		HOST:         hostName(),
		OS:           targetOS(),
		PID:          processID(),
		PROCESS:      processName(),
		SEVERITY:     ERROR_SEV,
		TIME:         timeNow(logger.timeFormat),
		FORWARD:      logger.errorForwardLogTo,
	}

	merged := meta.merge(fields...)
	str := merged.serialize()
	logger.write(str)
}

func (logger *FieldLogger) write(str string) {

	// Note: Making this faster is a good thing (while we are a sync logger - async logger is a different story)
	// So we don't use the stdlib logger.Print(), but rather have our own optimized version
	// Which does less, but is 3-10x faster

	// alloc a slice to contain the string and possible '\n'
	length := len(str)
	buffer := make([]byte, length+1)
	copy(buffer[:], str)
	if len(str) == 0 || str[length-1] != '\n' {
		copy(buffer[length:], "\n")
	}

	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	// This can return an error, but we just swallow it here as what can we or a client really do? Try and log it? :)
	logger.output.Write(buffer)
}
