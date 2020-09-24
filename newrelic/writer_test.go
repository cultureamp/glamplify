package newrelic

import (
	"context"
	"github.com/cultureamp/glamplify/log"
	"gotest.tools/assert"
	"os"
	"strings"
	"testing"
	"time"
)

func Test_NewRelic_Writer(t *testing.T) {
	if key := os.Getenv("NEW_RELIC_LICENSE_KEY"); key == "" {
		t.Skip("no new relic api key set")
	}

	ctx := context.Background()

	// https://log-api.newrelic.com/log/v1
	writer := newWriter(func(config *NRFieldWriter) {
		config.Endpoint ="https://log-api.newrelic.com/log/v1"
		config.Timeout =  time.Second * time.Duration(2)
	})
	logger := log.NewFromCtxWithCustomerWriter(ctx, writer)

	json := logger.Info("hello_world2")
	writer.WaitAll()

	assert.Assert(t, json != "", json)
	assert.Assert(t, strings.Contains(json, "hello"), json)
}
