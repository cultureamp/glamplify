package datadog

import (
	"context"
	gcontext "github.com/cultureamp/glamplify/context"
	"github.com/cultureamp/glamplify/log"
	"os"
	"sync"
	"testing"
	"time"
)

func Test_DataDog_Writer(t *testing.T) {

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

	writer := NewDataDogWriter(func(config *DDFieldWriter) {
		config.Endpoint = "https://http-intake.logs.datadoghq.com/v1/input"
		config.Timeout =  time.Second * time.Duration(2)
	})

	logger := log.NewFromCtxWithCustomerWriter(ctx, writer)

	logger.Info("hello Data Dog")
	writer.WaitAll()
}

func Test_DataDog_WaitGroup(t *testing.T) {

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

	writer := NewDataDogWriter(func(config *DDFieldWriter) {
		config.Endpoint = "https://http-intake.logs.datadoghq.com/v1/input"
		config.WaitGroup = &sync.WaitGroup{}
	})

	logger := log.NewFromCtxWithCustomerWriter(ctx, writer)

	logger.Info("hello Data Dog")
	writer.WaitAll()
}