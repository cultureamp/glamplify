package http

import (
	"gotest.tools/assert"
	"net/http"
	"testing"

	"github.com/cultureamp/glamplify/log"
	"github.com/cultureamp/glamplify/newrelic"
	"github.com/cultureamp/glamplify/bugsnag"
)

func Test_Wrap(t *testing.T) {

	app, appErr := newrelic.NewApplication("GlamplifyUnitTests", func(conf *newrelic.Config) {
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

	notifier, notifyErr := bugsnag.NewApplication("GlamplifyUnitTests", func (conf *bugsnag.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.AppVersion = "1.0.0"
	})
	assert.Assert(t, notifyErr == nil, notifyErr)

	pattern, handler := WrapHTTPHandler(app, notifier, "/", rootRequestHandler)
	assert.Assert(t, handler != nil, handler)
	assert.Assert(t, pattern == "/", pattern)

}

func rootRequestHandler(w http.ResponseWriter, r *http.Request) {}