package monitor

import (
	"errors"

	"github.com/cultureamp/glamplify/log"
)

// Logger is the interface that is used for logging in the go-agent.  Assign the
// Config.Logger field to the Logger you wish to use.  Loggers must be safe for
// use in multiple goroutines.  Two Logger implementations are included:
// NewLogger, which logs at info level, and NewDebugLogger which logs at debug
// level.  logrus and logxi are supported by the integration packages
// https://godoc.org/github.com/newrelic/go-agent/_integrations/nrlogrus and
// https://godoc.org/github.com/newrelic/go-agent/_integrations/nrlogxi/v1.
type monitorLogger struct {
	fieldLogger *log.FieldLogger
}

func newMonitorLogger() *monitorLogger {
	logger := log.New()

	return &monitorLogger{
		fieldLogger: logger,
	}
}

func (logger monitorLogger) Error(msg string, context map[string]interface{}) {
	err := errors.New(msg)
	logger.fieldLogger.Error(err, context)
}

func (logger monitorLogger) Warn(msg string, context map[string]interface{}) {
	logger.fieldLogger.Print(msg, context)
}

func (logger monitorLogger) Info(msg string, context map[string]interface{}) {
	logger.fieldLogger.Print(msg, context)
}

func (logger monitorLogger) Debug(msg string, context map[string]interface{}) {
	logger.fieldLogger.Debug(msg, context)
}

func (logger monitorLogger) DebugEnabled() bool {
	return false
}

func (logger monitorLogger) merge(logFields log.Fields, fields ...Fields) log.Fields {
	merged := log.Fields{}

	for k, v := range logFields {
		merged[k] = v
	}

	for _, f := range fields {
		for k, v := range f {
			merged[k] = v
		}
	}

	return merged
}