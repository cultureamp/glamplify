package bugsnag

import (
	"context"
	"fmt"

	"github.com/cultureamp/glamplify/log"
)

type bugsnagLogger struct {
	logger *log.Logger
}

func newBugsnagLogger(ctx context.Context) *bugsnagLogger {
	logger := log.NewFromCtx(ctx)

	return &bugsnagLogger{
		logger: logger,
	}
}

// Printf logs a message
func (l bugsnagLogger) Printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fields := log.Fields{
		log.Message: msg,
	}
	l.logger.Info("bugsnag_log", fields)
}
