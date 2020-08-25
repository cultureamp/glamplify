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

func (logger bugsnagLogger) Printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fields := log.Fields{
		log.Message: msg,
	}
	logger.logger.Info("bugsnag_log", fields)
}
