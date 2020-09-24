package datadog

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	gcontext "github.com/cultureamp/glamplify/context"
	"github.com/cultureamp/glamplify/log"
	"gotest.tools/assert"
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

	json := logger.Info("hello Data Dog")
	writer.WaitAll()

	assert.Assert(t, json != "", json)
	assert.Assert(t, strings.Contains(json, "hello"), json)
}