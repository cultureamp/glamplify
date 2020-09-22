package datadog

import (
	"context"
	"os"
	"testing"
	"time"

	gcontext "github.com/cultureamp/glamplify/context"
	"github.com/cultureamp/glamplify/log"
)

func Test_DataDog_RealWorld(t *testing.T) {

	if key := os.Getenv("DD_CLIENT_API_KEY"); key == "" {
		t.Skip("no data dog api key set")
	}

	ctx := context.Background()
	ctx = gcontext.AddRequestFields(ctx, gcontext.RequestScopedFields{
		TraceID:             "1-2-3",
		RequestID:           "7-8-9",
		CorrelationID:       "1-5-9",
		CustomerAggregateID: "hooli",
		UserAggregateID:     "UserAggregateID-123",
	})

	writer := NewDataDogWriter(func(config *writerConfig) {
		config.endpoint = "https://http-intake.logs.datadoghq.com/v1/input"
	})

	logger := log.NewFromCtxWithCustomerWriter(ctx, writer)

	logger.Info("hello Data Dog")
	time.Sleep(2 * time.Second)
}
