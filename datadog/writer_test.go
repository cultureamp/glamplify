package datadog

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	gcontext "github.com/cultureamp/glamplify/context"
	"github.com/cultureamp/glamplify/log"
	"github.com/stretchr/testify/assert"
)

func Test_DataDog_Writer(t *testing.T) {

	if key := os.Getenv(DDApiKey); key == "" {
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

	assert.NotEmpty(t, json)
	assert.True(t, strings.Contains(json, "hello"))
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