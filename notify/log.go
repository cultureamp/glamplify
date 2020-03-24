package notify

import (
	"context"
	"fmt"
	"github.com/cultureamp/glamplify/constants"

	"github.com/cultureamp/glamplify/log"
)

type notifyLogger struct {
	logger *log.Logger
}

func newNotifyLogger(ctx context.Context) *notifyLogger {
	logger := log.FromScope(ctx)

	return &notifyLogger{
		logger: logger,
	}
}

func (logger notifyLogger) Printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	fields := log.Fields{
		constants.MessageLogField : msg,
	}
	logger.logger.Info("notified", fields)
}
