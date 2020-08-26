package http

import (
	"context"
	"github.com/cultureamp/glamplify/sentry"
	"gotest.tools/assert"
	"net/http"
	"testing"

	"github.com/cultureamp/glamplify/log"
	"github.com/cultureamp/glamplify/newrelic"
	"github.com/cultureamp/glamplify/bugsnag"
)

func Test_Wrap_NR_Bugsnag(t *testing.T) {
	ctx := context.Background()

	app, appErr := newrelic.NewApplication(ctx, "GlamplifyUnitTests", func(conf *newrelic.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = false
		conf.Labels = newrelic.Labels{
			"asset":          log.Unknown,
			"classification": "restricted",
			"workload":       "development",
			"camp":           "amplify",
		}
	})
	assert.Assert(t, appErr == nil, appErr)

	bugsnag, berr := bugsnag.NewApplication(ctx, "GlamplifyUnitTests", func (conf *bugsnag.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.AppVersion = "1.0.0"
	})
	assert.Assert(t, berr == nil, berr)

	pattern, handler := WrapHTTPHandlerWithNewrelicAndBusgnag(app, bugsnag, "/", rootRequestHandler)
	assert.Assert(t, handler != nil, handler)
	assert.Assert(t, pattern == "/", pattern)

}

func Test_Wrap_NR_Sentry(t *testing.T) {
	ctx := context.Background()

	app, appErr := newrelic.NewApplication(ctx, "GlamplifyUnitTests", func(conf *newrelic.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = false
		conf.Labels = newrelic.Labels{
			"asset":          log.Unknown,
			"classification": "restricted",
			"workload":       "development",
			"camp":           "amplify",
		}
	})
	assert.Assert(t, appErr == nil, appErr)

	sentry, serr := sentry.NewApplication(ctx, "GlamplifyUnitTests", func (conf *sentry.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.AppVersion = "1.0.0"
	})
	assert.Assert(t, serr == nil, serr)

	pattern, handler := WrapHTTPHandlerWithNewrelicAndSentry(app, sentry, "/", rootRequestHandler)
	assert.Assert(t, handler != nil, handler)
	assert.Assert(t, pattern == "/", pattern)

}

func rootRequestHandler(w http.ResponseWriter, r *http.Request) {}