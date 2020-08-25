package newrelic

import (
	"testing"
	"time"

	"github.com/cultureamp/glamplify/log"
)

func Test_Realworld(t *testing.T) {
	// https://log-api.newrelic.com/log/v1
	writer := newWriter(func(config *writerConfig) {
		config.endpoint ="https://log-api.newrelic.com/log/v1"
	})
	mlog := log.NewFromCtxWithCustomerWriter(ctx, writer)

	mlog.Info("hello_world2")
	time.Sleep(2 * time.Second)
}

