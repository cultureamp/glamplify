package newrelic

import (
	"context"
	"errors"

	"github.com/cultureamp/glamplify/log"
)

// Logger is the interface that is used for logging in the New Relic go-agent.  Assign the
// config.Logger types to the Logger you wish to use.  Loggers must be safe for
// use in multiple goroutines.
type agentLogger struct {
	logger *log.Logger
}

func newAgentLogger(ctx context.Context) *agentLogger {
	logger := log.NewFromCtx(ctx)

	return &agentLogger{
		logger: logger,
	}
}

func (app agentLogger) Error(msg string, context map[string]interface{}) {
	err := errors.New(msg)
	app.logger.Error("monitor_error", err, context)
}

func (app agentLogger) Warn(msg string, context map[string]interface{}) {
	app.logger.Warn("monitor_warn", context, log.Fields{
		log.Message: msg,
	})
}

func (app agentLogger) Info(msg string, context map[string]interface{}) {
	app.logger.Info("monitor_info", context, context, log.Fields{
		log.Message: msg,
	})
}

func (app agentLogger) Debug(msg string, context map[string]interface{}) {
	app.logger.Debug("monitor_debug", context, context, log.Fields{
		log.Message: msg,
	})
}

func (app agentLogger) DebugEnabled() bool {
	return false
}