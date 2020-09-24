package newrelic

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cultureamp/glamplify/log"
	"gotest.tools/assert"
)

func Test_NewRelic_Writer(t *testing.T) {
	if key := os.Getenv("NEW_RELIC_LICENSE_KEY"); key == "" {
		t.Skip("no new relic api key set")
	}

	ctx := context.Background()

	// https://log-api.newrelic.com/log/v1
	writer := NewNRWriter(func(config *NRFieldWriter) {
		config.Endpoint ="https://log-api.newrelic.com/log/v1"
		config.Timeout =  time.Second * time.Duration(2)
	})
	logger := log.NewFromCtxWithCustomerWriter(ctx, writer)

	json := logger.Info("hello_world2")
	writer.WaitAll()

	assert.Assert(t, json != "", json)
	assert.Assert(t, strings.Contains(json, "hello"), json)
}

func Test_NewRelic_Writer_IsEnabled(t *testing.T) {

	writer := NewNRWriter(func(config *NRFieldWriter) {
		config.Level = log.WarnSev
	})

	ok := writer.IsEnabled(log.DebugSev)
	assert.Assert(t, !ok, ok)
	ok = writer.IsEnabled(log.InfoSev)
	assert.Assert(t, !ok, ok)
	ok = writer.IsEnabled(log.WarnSev)
	assert.Assert(t, ok, ok)
	ok = writer.IsEnabled(log.ErrorSev)
	assert.Assert(t, ok, ok)
	ok = writer.IsEnabled(log.FatalSev)
	assert.Assert(t, ok, ok)
}