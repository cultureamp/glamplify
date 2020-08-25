package http

import (
	"github.com/cultureamp/glamplify/newrelic"
	"github.com/cultureamp/glamplify/bugsnag"
	"net/http"
)

func WrapHTTPHandler(
	app *newrelic.Application,
	notify *bugsnag.Application,
	pattern string,
	handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {

	// 1. Wrap with bugsnag
	pattern, handler = notify.WrapHTTPHandler(pattern, handler)

	// 2. Then wrap with new relic
	return app.WrapHTTPHandler(pattern, handler)
}
