package newrelic

import (
	"context"
	"github.com/cultureamp/glamplify/log"
	"os"
	"sync"
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
	mlog := log.NewFromCtxWithCustomerWriter(ctx, writer)

	mlog.Info("hello_world2")
	writer.WaitAll()
}

func Test_NewRelic_WaitGroup(t *testing.T) {
	if key := os.Getenv("NEW_RELIC_LICENSE_KEY"); key == "" {
		t.Skip("no new relic api key set")
	}

	ctx := context.Background()

	// https://log-api.newrelic.com/log/v1
	writer := newWriter(func(config *NRFieldWriter) {
		config.Endpoint ="https://log-api.newrelic.com/log/v1"
		config.WaitGroup = &sync.WaitGroup{}
	})
	mlog := log.NewFromCtxWithCustomerWriter(ctx, writer)

	mlog.Info("hello_world2")

	writer.WaitAll()
}
