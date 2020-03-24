package monitor

import (
	"context"
	"errors"

	"github.com/cultureamp/glamplify/log"
)

// Logger is the interface that is used for logging in the New Relic go-agent.  Assign the
// config.Logger types to the Logger you wish to use.  Loggers must be safe for
// use in multiple goroutines.  Two Logger implementations are included:
// NewLogger, which logs at info level, and NewDebugLogger which logs at debug
// level.  logrus and logxi are supported by the integration packages
// https://godoc.org/github.com/newrelic/go-agent/_integrations/nrlogrus and
// https://godoc.org/github.com/newrelic/go-agent/_integrations/nrlogxi/v1.
type monitorLogger struct {
	scope *log.Logger
}

func newMonitorLogger(ctx context.Context) *monitorLogger {
	scope := log.FromScope(ctx)

	return &monitorLogger{
		scope: scope,
	}
}

func (logger monitorLogger) Error(msg string, context map[string]interface{}) {
	err := errors.New(msg)
	logger.scope.Error(err, context)
}

func (logger monitorLogger) Warn(msg string, context map[string]interface{}) {
	logger.scope.Warn(msg, context)
}

func (logger monitorLogger) Info(msg string, context map[string]interface{}) {
	logger.scope.Info(msg, context)
}

func (logger monitorLogger) Debug(msg string, context map[string]interface{}) {
	logger.scope.Debug(msg, context)
}

func (logger monitorLogger) DebugEnabled() bool {
	return false
}