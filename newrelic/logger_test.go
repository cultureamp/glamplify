package newrelic

import (
	"context"
	"errors"
	gcontext "github.com/cultureamp/glamplify/context"
	"github.com/cultureamp/glamplify/log"
	"os"
	"testing"
)

var (
	ctx      context.Context
	rsFields gcontext.RequestScopedFields
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	ctx = context.Background()
	ctx = gcontext.AddRequestFields(ctx, gcontext.RequestScopedFields{
		TraceID:             "1-2-3",
		RequestID:           "7-8-9",
		CorrelationID:       "1-5-9",
		CustomerAggregateID: "hooli",
		UserAggregateID:     "UserAggregateID-123",
	})

	rsFields, _ = gcontext.GetRequestScopedFields(ctx)

	os.Setenv("PRODUCT", "engagement")
	os.Setenv("APP", "murmur")
	os.Setenv("APP_VERSION", "87.23.11")
	os.Setenv("AWS_REGION", "us-west-02")
	os.Setenv("AWS_ACCOUNT_ID", "aws-account-123")
}

func shutdown() {
	os.Unsetenv("PRODUCT")
	os.Unsetenv("APP")
	os.Unsetenv("APP_VERSION")
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("AWS_ACCOUNT_ID")
}


func Test_Monitor_Debug(t *testing.T) {

	logger := NewFromCtx(ctx)
	logger.Event("monitor_debug").Fields(log.Fields{"name": "foo"}).Debug("debug message")
}

func Test_Monitor_Info(t *testing.T) {
	logger := NewFromCtx(ctx)
	logger.Event("monitor_info").Fields(log.Fields{"name": "foo"}).Info("info message")
}

func Test_Monitor_Warn(t *testing.T) {
	logger := NewFromCtx(ctx)
	logger.Event("monitorwarn").Fields(log.Fields{"name": "foo"}).Warn("warn message")
}

func Test_Monitor_Error(t *testing.T) {
	logger := NewFromCtx(ctx)
	logger.Event("monitor_error").Fields(log.Fields{"name": "foo"}).Error(errors.New("error message"))
}

func Test_Monitor_Fatal(t *testing.T) {

	defer func() {
		recover()
	}()

	logger := NewFromCtx(ctx)
	logger.Event("monitor_fatal").Fields(log.Fields{"name": "foo"}).Fatal(errors.New("fatal message"))
}


