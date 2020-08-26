package http

import (
	"github.com/cultureamp/glamplify/newrelic"
	"github.com/cultureamp/glamplify/bugsnag"
	"github.com/cultureamp/glamplify/sentry"
	"net/http"
)

func WrapHTTPHandlerWithNewrelicAndBusgnag(
	app *newrelic.Application,
	bugsnap *bugsnag.Application,
	pattern string,
	handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {

	// 1. Wrap with bugsnag
	pattern, handler = bugsnap.WrapHTTPHandler(pattern, handler)

	// 2. Then wrap with new relic
	return app.WrapHTTPHandler(pattern, handler)
}

func WrapHTTPHandlerWithNewrelicAndSentry(
	app *newrelic.Application,
	sentry *sentry.Application,
	pattern string,
	handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {

	// 1. Wrap with sentry
	pattern, handler = sentry.WrapHTTPHandler(pattern, handler)

	// 2. Then wrap with new relic
	return app.WrapHTTPHandler(pattern, handler)
}
