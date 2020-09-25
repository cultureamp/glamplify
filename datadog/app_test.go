package datadog_test

import (
	"context"
	"testing"

	"github.com/cultureamp/glamplify/datadog"
	"gotest.tools/assert"
)

func Test_DataDog_Application(t *testing.T) {
	ctx := context.Background()
	app, err := datadog.NewApplication(ctx, "Glamplify-Unit-Tests", func(conf *datadog.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = false
	})

	assert.Assert(t, err == nil, err)
	assert.Assert(t, app != nil, "application was nil")

	app.Shutdown()
}

