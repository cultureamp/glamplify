package datadog

import (
	"context"
	"github.com/cultureamp/glamplify/env"
	"os"
	"testing"
	"time"

	gcontext "github.com/cultureamp/glamplify/context"
	"github.com/cultureamp/glamplify/log"
	"github.com/stretchr/testify/assert"
)

func Test_DataDog_Writer(t *testing.T) {

	if key := os.Getenv(env.DatadogApiKey); key == "" {
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

	json := logger.Event("data_dog").Fields(log.Fields{
		"string_key": "a string",
		"int_key": 123,
	}).Info("hello from Data Dog")
	writer.WaitAll()

	assert.NotEmpty(t, json)
	assert.Contains(t, json, "\"event\":\"data_dog\"")
	assert.Contains(t, json, "\"string_key\":\"a string\"")
	assert.Contains(t, json, "\"int_key\":123")
	assert.Contains(t, json, "hello from Data Dog")
}

func Test_DataDog_Writer_IsEnabled(t *testing.T) {

	writer := NewDataDogWriter(func(config *DDFieldWriter) {
		config.Level = log.ErrorSev
	})

	ok := writer.IsEnabled(log.DebugSev)
	assert.False(t, ok)
	ok = writer.IsEnabled(log.InfoSev)
	assert.False(t, ok)
	ok = writer.IsEnabled(log.WarnSev)
	assert.False(t, ok)
	ok = writer.IsEnabled(log.ErrorSev)
	assert.True(t, ok)
	ok = writer.IsEnabled(log.FatalSev)
	assert.True(t, ok)
}