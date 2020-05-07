package main

import (
	"errors"
	"github.com/cultureamp/glamplify/aws"
	"github.com/cultureamp/glamplify/config"
	http2 "github.com/cultureamp/glamplify/http"
	"github.com/cultureamp/glamplify/jwt"
	"github.com/cultureamp/glamplify/log"
	"github.com/cultureamp/glamplify/monitor"
	"github.com/cultureamp/glamplify/notify"
	"net/http"
	"net/http/httptest"
)

func main() {

	/* CONFIG */

	// settings will contain configuration data as read in from the config file.
	settings := config.Load()

	// Or if you want to look for a config file from a specific location use
	//settings = config.LoadFrom([]string{"${HOME}/settings"}, "config")

	// Then you can use
	if settings.App.Version > 2.0 {
		// to do
	}

	/* LOGGING */
	// Creating loggers is cheap. Create them on every request/run
	// DO NOT CACHE/REUSE THEM
	transactionFields := log.RequestScopedFields{
		TraceID:             "abc",   // Get TraceID from context or from wherever you have it stored
		UserAggregateID:     "user1", // Get UserAggregateID from context or from wherever you have it stored
		CustomerAggregateID: "cust1", // Get CustomerAggregateID from context or from wherever you have it stored
	}
	logger := log.New(transactionFields)

	// or if you want a field to be present on each subsequent logging call do this:
	logger = log.New(transactionFields, log.Fields{"request_id": 123})

	/* Monitor & Notify */
	app, appErr := monitor.NewApplication("GlamplifyUnitTests", func(conf *monitor.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = false
		conf.Labels = monitor.Labels{
			"asset":          log.Unknown,
			"classification": "restricted",
			"workload":       "development",
			"camp":           "amplify",
		}
	})
	if appErr != nil {
		logger.Fatal(appErr)
	}

	notifier, notifyErr := notify.NewNotifier("GlamplifyUnitTests", func (conf *notify.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.AppVersion = "1.0.0"
	})
	if notifyErr != nil {
		logger.Fatal(notifyErr)
	}

	pattern, handler := http2.WrapHTTPHandler(app, notifier, "/", rootRequestHandler)
	h := http.HandlerFunc(handler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", pattern, nil)
	h.ServeHTTP(rr, req)

	app.Shutdown()
}

func rootRequestHandler(w http.ResponseWriter, r *http.Request) {

	// get JWT payload from http header
	payload, err := jwt.PayloadFromRequest(r)

	// Create the logging config for this request
	ctx := r.Context()
	traceID, _ := aws.GetTraceID(ctx)
	requestScopedFields := log.RequestScopedFields{
		TraceID:             traceID,				// Get TraceID from context or from wherever you have it stored
		UserAggregateID:     payload.EffectiveUser, // Get UserAggregateID from context or from wherever you have it stored
		CustomerAggregateID: payload.Customer,      // Get CustomerAggregateID from context or from wherever you have it stored
	}

	// Then create a logger that will use those transaction fields values when writing out logs
	logger := log.New(requestScopedFields)

	// OR if you want a helper to do all of the above, use
	r = log.EnsureRequestScopedFieldsPresentInRequest(r)
	logger = log.NewFromRequest(r)

	logger.Debug("something_happened")

	// optional: save this to the context for later use, then you can just create via: ctx, logger := log.NewFromCtx(ctx)
	// ctx = log.AddRequestScopedFieldsCtx(ctx, requestScopedFields)
	// if you want to get them back out from the context
	// rsFields := log.GetRequestScopedFieldsCtx(ctx)

	// optional: if you need to propagate the request then make sure you update the context for the request
	// then you can create new loggers with: r, logger := log.NewFromRequest(r)
	// r = r.WithContext(ctx)


	// or use the default logger with transaction fields passed in
	log.Debug(requestScopedFields, "something_happened", log.Fields{})

	// Emit debug trace with types
	// Fields can contain any type of variables
	logger.Debug("something_else_happened", log.Fields{
		"aString": "hello",
		"aInt":    123,
		"aFloat":  42.48,
	})

	// Emit normal logging (can add optional types if required)
	// Typically Print will be sent onto 3rd party aggregation tools (eg. Splunk)
	logger.Info("Executing main")

	// Emit Error (can add optional types if required)
	// Errors will always be sent onto 3rd party aggregation tools (eg. Splunk)
	err = errors.New("failed to save record to db")
	logger.Error(err)

	// Emit Fatal (can add optional types if required) and PANIC!
	// Fatal error will always be sent onto 3rd party aggregation tools (eg. Splunk)
	//err = errors.New("program died")
	//logger.Fatal(err)

	/* NEW RELIC TRANSACTION */
	txn, err := monitor.TxnFromRequest(w, r)
	if err == nil {
		txn.AddAttributes(log.Fields{
			"aString": "hello world",
			"aInt":    123,
		})
	}

	// Do more things

	/* NEW RELIC Add Attributes */
	if err == nil {
		txn.AddAttributes(log.Fields{
			"aString2": "goodbye",
			"aInt2":    456,
		})
	}
}
