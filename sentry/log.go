package bugsnag

import (
	"context"
	"fmt"

	"github.com/cultureamp/glamplify/log"
)

type notifyLogger struct {
	logger *log.Logger
}

func newNotifyLogger(ctx context.Context) *notifyLogger {
	logger := log.NewFromCtx(ctx)

	return &notifyLogger{
		logger: logger,
	}
}

func (logger notifyLogger) Printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fields := log.Fields{
		log.Message: msg,
	}
	logger.logger.Info("notified", fields)
}
