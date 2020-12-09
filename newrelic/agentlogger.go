package newrelic

import (
	"context"
	"errors"

	"github.com/cultureamp/glamplify/log"
)

type agentLogger struct {
	logger *log.Logger
}

func newAgentLogger(ctx context.Context) *agentLogger {
	logger := log.NewFromCtx(ctx)

	return &agentLogger{
		logger: logger,
	}
}

// Error logs an error
func (app agentLogger) Error(msg string, context map[string]interface{}) {
	err := errors.New(msg)
	app.logger.Error("monitor_error", err, context)
}

// Warn logs a warning
func (app agentLogger) Warn(msg string, context map[string]interface{}) {
	app.logger.Warn("monitor_warn", context, log.Fields{
		log.Message: msg,
	})
}

// Info logs an info message
func (app agentLogger) Info(msg string, context map[string]interface{}) {
	app.logger.Info("monitor_info", context, context, log.Fields{
		log.Message: msg,
	})
}

// Debug logs a debug message
func (app agentLogger) Debug(msg string, context map[string]interface{}) {
	app.logger.Debug("monitor_debug", context, context, log.Fields{
		log.Message: msg,
	})
}

// DebugEnabled returns true if debugging is enabled, false otherwise
func (app agentLogger) DebugEnabled() bool {
	return false
}