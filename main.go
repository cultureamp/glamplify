package main

import (
	"context"
	"errors"
	"github.com/cultureamp/glamplify/datadog"
	"github.com/cultureamp/glamplify/sentry"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/cultureamp/glamplify/aws"
	gcontext "github.com/cultureamp/glamplify/context"
	"github.com/cultureamp/glamplify/jwt"
	"github.com/cultureamp/glamplify/log"
	"github.com/cultureamp/glamplify/newrelic"
)

func main() {
	ctx := context.Background()

	/****** XRAY ******/
	xrayTracer := aws.NewTracer(ctx, func(conf *aws.TracerConfig) {
		conf.Environment = "production" // or "development"
		conf.AWSService = "ECS"         // or "EC2" or "LAMBDA"
		conf.EnableLogging = true
		conf.Version = os.Getenv("APP_VERSION")
	})

	/****** LOGGING ******/

	// Creating loggers is cheap. Create them on every request/run
	// DO NOT CACHE/REUSE THEM
	transactionFields := gcontext.RequestScopedFields{
		TraceID:             "abc",   // Get TraceID from AWS Xray
		RequestID:			 "req1", // Get RequestID from X-Request-ID header
		CorrelationID:		 "uuid4", // Get CorrelationID from X-Correlation-ID header (web-gateway will add this in)
		UserAggregateID:     "user1", // Get UserAggregateID from context or from wherever you have it stored
		CustomerAggregateID: "cust1", // Get CustomerAggregateID from context or from wherever you have it stored
	}
	logger := log.New(transactionFields)

	// or if you want a field to be present on each subsequent logging call do this:
	logger = log.New(transactionFields, log.Fields{"request_id": 123})

	/* DataDog */
	datadogApp := datadog.NewApplication(ctx, "GlamplifyUnitTests", func(conf *datadog.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = false
	})

	sentryApp, sentryErr := sentry.NewApplication(ctx, "GlamplifyUnitTests", func(conf *sentry.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.AppVersion = "1.0.0"
	})
	if sentryErr != nil {
		logger.Fatal("sentry_failed", sentryErr)
	}

	h := xrayTracer.SegmentHandler("MyApp", sentryApp.Middleware(datadogApp.WrapHandler("MyApp", rootRequestHandler)))

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(rr, req)

	datadogApp.Shutdown()
	sentryApp.Shutdown()
}

func rootRequestHandler(w http.ResponseWriter, r *http.Request) {

	// get JWT payload from http header
	decoder, err := jwt.NewDecoder() // assumes AUTH_PUBLIC_KEY set, check other New methods for overloads
	if err != nil {
		// handle error
	}
	payload, err := jwt.PayloadFromRequest(r, decoder)
	if err != nil {
		// handle error
	}
	
	// Create the logging config for this request
	requestScopedFields := gcontext.RequestScopedFields{
		TraceID:             r.Header.Get(gcontext.TraceIDHeader),
		RequestID:           r.Header.Get(gcontext.RequestIDHeader),
		CorrelationID:       r.Header.Get(gcontext.CorrelationIDHeader),
		UserAggregateID:     payload.EffectiveUser,
		CustomerAggregateID: payload.Customer,
	}
	logger := log.New(requestScopedFields) // Then create a logger that will use those transaction fields values when writing out logs

	// OR if you want a helper to do all of the above, use
	r = gcontext.WrapRequest(r)
	logger = log.NewFromRequest(r)

	// now away you go!
	logger.Debug("something_happened")

	// or use the default logger with transaction fields passed in
	log.Debug(requestScopedFields, "something_happened", log.Fields{log.Message: "message"})

	// or use an more expressive syntax
	logger.Event("something_happened").Debug("message")

	// or use an more expressive syntax
	logger.Event("something_happened").Fields(log.Fields{"count": 1}).Debug("message")

	// Emit debug trace with types
	// Fields can contain any type of variables
	logger.Debug("something_else_happened", log.Fields{
		"aString":   "hello",
		"aInt":      123,
		"aFloat":    42.48,
		log.Message: "message",
	})
	logger.Event("something_else_happened").Fields(log.Fields{
		"aString": "hello",
		"aInt":    123,
		"aFloat":  42.48,
	}).Debug("message")

	// Emit normal logging (can add optional types if required)
	// Typically Print will be sent onto 3rd party aggregation tools (eg. Splunk)
	logger.Info("Executing main")

	// Emit Error (can add optional types if required)
	// Errors will always be sent onto 3rd party aggregation tools (eg. Splunk)
	err = errors.New("failed to save record to db")
	logger.Error("error_saving_record", err)

	// Emit Fatal (can add optional types if required) and PANIC!
	// Fatal error will always be sent onto 3rd party aggregation tools (eg. Splunk)
	//err = errors.New("program died")
	//logger.Fatal("fatal_problem_detected", err)

	/* NEW RELIC TRANSACTION */
	txn, err := newrelic.TxnFromRequest(w, r)
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
