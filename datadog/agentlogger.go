package datadog

import (
	"context"
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

func (log agentLogger) Log(msg string) {

}
