package newrelic

import (
	"context"
	"errors"

	"github.com/cultureamp/glamplify/log"
)

// Logger is the interface that is used for logging in the New Relic go-agent.  Assign the
// config.Logger types to the Logger you wish to use.  Loggers must be safe for
// use in multiple goroutines.
type monitorLogger struct {
	logger *log.Logger
}

func newMonitorLogger(ctx context.Context) *monitorLogger {
	logger := log.NewFromCtx(ctx)

	return &monitorLogger{
		logger: logger,
	}
}

func (app monitorLogger) Error(msg string, context map[string]interface{}) {
	err := errors.New(msg)
	app.logger.Error("monitor_error", err, context)
}

func (app monitorLogger) Warn(msg string, context map[string]interface{}) {
	app.logger.Warn("monitor_warn", context, log.Fields{
		log.Message: msg,
	})
}

func (app monitorLogger) Info(msg string, context map[string]interface{}) {
	app.logger.Info("monitor_info", context, context, log.Fields{
		log.Message: msg,
	})
}

func (app monitorLogger) Debug(msg string, context map[string]interface{}) {
	app.logger.Debug("monitor_debug", context, context, log.Fields{
		log.Message: msg,
	})
}

func (app monitorLogger) DebugEnabled() bool {
	return false
}