package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"

	awsxray "github.com/aws/aws-xray-sdk-go/xray"
	"github.com/cultureamp/glamplify/aws"
	gcontext "github.com/cultureamp/glamplify/context"
	"github.com/cultureamp/glamplify/datadog"
	"github.com/cultureamp/glamplify/jwt"
	"github.com/cultureamp/glamplify/log"
	"github.com/cultureamp/glamplify/sentry"
	"github.com/google/uuid"
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
		panic(sentryErr)
	}

	h := xrayTracer.SegmentHandler("MyApp", sentryApp.Middleware(datadogApp.WrapHandler("MyApp", rootRequestHandler)))

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(rr, req)

	datadogApp.Shutdown()
	sentryApp.Shutdown()
}

func rootRequestHandler(w http.ResponseWriter, r *http.Request) {
	ctx := GetLoggingCtx(r.Context(), r)
	logger := log.NewFromCtx(ctx)

	// get JWT payload from http header
	decoder, err := jwt.NewDecoder() // assumes AUTH_PUBLIC_KEY set, check other New methods for overloads
	if err != nil {
		logger.Event("glamplify_request_handler").Error(err)
	}
	payload, err := jwt.PayloadFromRequest(r, decoder)
	if err != nil {
		logger.Event("glamplify_request_handler").Error(err)
	}

	// Fields can contain any type of variables
	logger.Event("glamplify_request_handler").Fields(log.Fields{
		"payload": payload,
		"aString":   "hello",
		"aInt":      123,
		"aFloat":    42.48,
	}).Debug("payload")
}

// GetLoggingCtx adds in missing TraceID, RequestID and CorrelationID if required at the start of a request
func GetLoggingCtx(ctx context.Context, r *http.Request) context.Context {

	// If XRAY is enabled, get the trace_id
	traceID := getTraceID(ctx, r)

	// Get the customer RequestID or set to UUID if empty
	requestID := getRequestID(r)

	// Set to UUID for internal usage
	correlationID := getCorrelationID()

	rsFields := gcontext.RequestScopedFields{
		TraceID: traceID,
		RequestID: requestID,
		CorrelationID: correlationID,
	}

	return gcontext.AddRequestFields(ctx, rsFields)
}


func getTraceID(ctx context.Context,  r *http.Request) string {

	if awsxray.RequestWasTraced(ctx) {
		return awsxray.TraceID(ctx)
	}

	traceID := r.Header.Get(gcontext.TraceIDHeader)
	if traceID == "" {
		return gcontext.ErrorUUID
	}

	return traceID
}

func getRequestID( r *http.Request) string {

	// Did the client pass us a request_id?
	requestID := r.Header.Get(gcontext.RequestIDHeader)
	if requestID == "" {
		// If client has set then we honour that value, otherwise we set it to a UUID
		requestID = uuid.New().String()
	}

	return requestID
}

func getCorrelationID() string {
	return uuid.New().String()
}
