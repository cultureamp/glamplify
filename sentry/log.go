package sentry

import (
	"context"
	"github.com/cultureamp/glamplify/log"
)

type sentryLogger struct {
	logger *log.Logger
}

func newSentryLogger(ctx context.Context) *sentryLogger {
	logger := log.NewFromCtx(ctx)

	return &sentryLogger{
		logger: logger,
	}
}

func (l sentryLogger) Write(p []byte) (n int, err error) {
	msg := string(p)
	l.logger.Event("sentry_log").Debug(msg)

	return len(p), nil
}

