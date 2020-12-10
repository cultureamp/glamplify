# glamplify
Go Amplify Module of useful common tools. The guiding principle is to implement a very light weight wrapper over the standard library (or where not adequate an open source community library), that conforms to our standard practises (12-Factor) and sensible defaults.


## Install

```
go get github.com/cultureamp/glamplify
```

## Usage

### TRACER (AWS XRAY)

```GO
package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"

    "github.com/aws/aws-xray-sdk-go/xray"
	"github.com/cultureamp/glamplify/aws"
	"github.com/cultureamp/glamplify/config"
	gcontext "github.com/cultureamp/glamplify/context"
	ghttp "github.com/cultureamp/glamplify/http"
	"github.com/cultureamp/glamplify/jwt"
	"github.com/cultureamp/glamplify/log"
	"github.com/cultureamp/glamplify/newrelic"
	"github.com/cultureamp/glamplify/bugsnag"
)

func main() {
	ctx := context.Background()

	xrayTracer := aws.NewTracer(ctx, func(conf *aws.TracerConfig) {
		conf.Environment = "production" // or "development"
		conf.AWSService = "ECS"         // or "EC2" 
		conf.EnableLogging = true
		conf.Version = os.Getenv("APP_VERSION")
	})
	
    h := xrayTracer.SegmentHandler("MyApp", http.HandlerFunc(requestHandler))

    if err := http.ListenAndServe(":8080", h); err != nil {
        panic(err)
    }
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
    // DO STUFF
    ctx := r.Context()
    
    // If you want to trace a critical section of cod, use
     xray.Capture(ctx, "segment-name", func(ctx1 context.Context) error {
    
       // DO THE THINGS
        var result interface{}
    
        return xray.AddMetadata(ctx1, "ResourceResult", result)
      })

    // DO MORE STUFF

}

```

### Logging

Logging in GO supports the Culture Amp [sensible default](https://cultureamp.atlassian.net/wiki/spaces/TV/pages/959939199/Logging)

```Go
package main

import (
    "bytes"
    "context"
    "errors"
    "net/http"
    "time"

    "github.com/cultureamp/glamplify/aws"
    "github.com/cultureamp/glamplify/constants"
    gcontext "github.com/cultureamp/glamplify/context"
    "github.com/cultureamp/glamplify/jwt"
    "github.com/cultureamp/glamplify/log"
)

func main() {

    // Creating loggers is cheap. Create them on every request/run
    // DO NOT CACHE/REUSE THEM
    transactionFields := gcontext.RequestScopedFields{
        TraceID:                "abc",          // Get TraceID from AWS xray 
        RequestID:              "random-string",// X-Request-ID, set optionally by clients
        CorrelationID:          "uuid4",        // X-Correlation-ID set by web-gateway as UUID4
        UserAggregateID :       "user1",        // Get User from JWT 
        CustomerAggregateID:    "cust1",        // Get Customer from JWT
   	}
    logger := log.New(transactionFields)

    // Or if you want a field to be present on each subsequent logging call do this:
    logger = log.New(transactionFields, log.Fields{"request_id": 123})

    h := http.HandlerFunc(requestHandler)

    if err := http.ListenAndServe(":8080", h); err != nil {
        logger.Error("failed_to_serve_http_request", err)
    }
}

func requestHandler(w http.ResponseWriter, r *http.Request) {

    // get JWT payload from http header
    decoder, err := jwt.NewDecoder()
    payload, err := jwt.PayloadFromRequest(r, decoder)
    
    // Create the logging config for this request
    requestScopedFields := gcontext.RequestScopedFields{
        TraceID: r.Header.Get(gcontext.TraceIDHeader),				
        RequestID: r.Header.Get(gcontext.RequestIDHeader),
        CorrelationID: r.Header.Get(gcontext.CorrelationIDHeader),      
        UserAggregateID: payload.EffectiveUser,     
        CustomerAggregateID: payload.Customer,      
    }
    
    // Then create a logger that will use those transaction fields values when writing out logs
    logger := log.New(requestScopedFields)

    // OR, if you want a helper that does all of the above, use
    r = gcontext.WrapRequest(r)  // Reads and Sets TraceID, RequestID, CorrelationID, User and Customer
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
        "aString": "hello",
        "aInt":    123,
        "aFloat":  42.48,
        log.Message: "message",
    })
    logger.Event("something_else_happened").Fields(log.Fields{
        "aString": "hello",
        "aInt":    123,
        "aFloat":  42.48,
    }).Debug("message")

    // Fields can contain any type of variables, but here are some helpful predefined ones
    // (see constants.go for full list)
    // MessageLogField             = "message"
    // TimeTakenLogField           = "time_taken"
    // MemoryUsedLogField          = "memory_used"
    // MemoryAvailLogField         = "memory_available"
    // ItemsProcessedLogField      = "items_processed"
    // TotalItemsProcessedLogField = "total_items_processed"
    // TotalItemsRequestedLogField = "total_items_requested"

    d := time.Millisecond * 123
    logger.Event("something_happened").Fields(log.Fields{
        log.TimeTaken : log.DurationAsISO8601(d), // returns "P0.123S" as per sensible default
        log.User: "MMLKSN443FN",
        "report":  "NVJKSJFJ34NBFN44",
        "aInt":    123,
        "aFloat":  42.48,
        "aString": "more info",
     }).Debug("The thing did what we expected it to do")

    // Typically Info will be sent onto 3rd party aggregation tools (eg. Splunk)
    logger.Info("something_happened_event")

    // Fields can contain any type of variables
    d = time.Millisecond * 456
    logger.Info("Something_happened_event", log.Fields{
        "program-name": "helloworld.exe",
        "start-up-param":    123,
        log.User:  "admin",
        log.Message: "The thing did what we expected it to do",
        log.TimeTaken: log.DurationAsISO8601(d), // returns "P0.456S"
    })

    // Errors will always be sent onto 3rd party aggregation tools (eg. Splunk)
    err = errors.New("missing database connection string")
    logger.Error("database_connection", err)

    // Fields can contain any type of variables
    err = errors.New("missing database connection string")
    logger.Event("database_connection").Fields(log.Fields{
        "program-name": "helloworld.exe",
        "start-up-param":    123,
        log.User:  "admin",
        log.Message: "The thing did not do what we expected it to do",
     }).Error(err)
}

```
Use `log.New()` for logging without a http request. Use `log.NewFromRequest()` when you do have a http request. This initializes a bunch of stuff for you (eg. JWT details are automatically logged for you). NewFromRequest always returns a valid Logger and an optional error (which usually describes problems decoding the JWT etc)

Use `Debug` for logging that will only be used when diving deep to uncover bugs. Typically `scope.Debug` messages will not automatically be sent to other 3rd party systems (eg. Splunk).

Use `Info` for standard log messages that you want to see always. These will never be turned off and will likely be always sent to 3rd party systems for further analysis (eg. Spliunk).

Use `Warn` when you have encounter a warning that should be looked at by a human, but has been recovered. All warning messages will be forwarded to 3rd party systems for monitoring and further analysis.

Use `Error` when you have encountered a GO error. This will NOT stop the program, it is assumed that the system has recovered. All error messages will be forwarded to 3rd party systems for monitoring and further analysis.

Use `Fatal` when you have encountered a GO error that is not recoverable. This will stop the program by calling panic(). All fatal messages will be forwarded to 3rd party systems for monitoring and further analysis.

use 'Audit' when you want to publish this log to our external customer facing API where they can retrieve information about what is happening in their account. 

### Lambda

```go

import (
    aws "github.com/aws/aws-lambda-go/events"
	gaws "github.com/cultureamp/glamplify/aws"
    "github.com/cultureamp/glamplify/datadog"
    "github.com/cultureamp/glamplify/log"
    "github.com/cultureamp/glamplify/sentry"
)

func main() {
    xray = gaws.NewTracer(ctx, func(config *gaws.TracerConfig) {
        config.EnableLogging = false
        config.Version = "1.0.0"
        config.Environment = "production"
        config.AWSService = "ECS"
    })

    // https://docs.sentry.io/platforms/go/serverless/ 
    //Sentry doesn't write to CloudWatch like other tools (eg. DataDog) 
    //So we have fo make sure we flush any pending data to Sentry before the lambda is completed and AWS close the network on us 
    //This is really yuk, but we either have to sentry.Flush() at the end of every handler, or use a HttpSyncTransport 
    //I'm not sure which is best... :(
    sentry, err := sentry.NewApplication(ctx, settings.App, func(config *sentry.Config) {
        config.Enabled = true
        config.Logging = false
        config.DSN = os.Getenv("SENTRY_DSN")
        config.Transport =  &sentrygo.HTTPSyncTransport{Timeout: 100 * time.Millisecond}
    })
    
    datadog = datadog.NewApplication(ctx, settings.App, func(conf *datadog.Config) {
        conf.Enabled = true
        conf.Logging = false
        conf.APIKey =  os.Getenv("DD_API_KEY")
        conf.WithRuntimeMetrics = true
        conf.Tags = datadog.Tags{"app": "myapp-api"}
        conf.ServerlessMode = true
    })
    
    lambda.Start(datadogApp.WrapLambdaHandler(handler))
}

func handler(ctx context.Context, request aws.ALBTargetGroupRequest) (aws.ALBTargetGroupResponse, error) {
    // https://docs.sentry.io/platforms/go/serverless/
    // Make sure we catch any panics and report them to Sentry...
    defer sentrygo.Recover()
    
    ctx = getLoggingCtx(ctx, request)
    span, ctx := datadogApp.TraceHandler(ctx, request.Path)
    defer span.Finish()
    
    logger := log.NewFromCtx(ctx)
    logger.Event("myapp-api").Fields(log.Fields{
        "path":        request.Path,
        "http_method": request.HTTPMethod,
    }).Debug("Starting request...")
    
    var response aws.ALBTargetGroupResponse
    var err error
    
    start := time.Now()
    switch request.HTTPMethod {
    case "GET":
        response, err = handleGet(ctx, request)
    default:
        response, err = unhandled(ctx, request)
    }
    duration := time.Since(start)

    // Returning nil Headers and/or MultiValueHeaders causing the ALB to reject the response
    // So patch the response with default values if any mandatory fields are missing
    response = patch(response)

    logger.Event("authz_api").Fields(log.Fields{
        "status":             response.StatusCode,
        "status_description": response.StatusDescription,
        "is_base64_encoded":  response.IsBase64Encoded,
        "body":               response.Body,
    }).Fields(log.NewDurationFields(duration)).Debug("Finished request")
    
    return response, nil
}

func patch(response aws.ALBTargetGroupResponse) aws.ALBTargetGroupResponse {
    // https://serverless-training.com/articles/api-gateway-vs-application-load-balancer-technical-details/
    
    // Returning nil Headers and/or MultiValueHeaders causing the ALB to reject the response
    // So make sure they are set if nil to empty maps...
    if response.Headers == nil {
        response.Headers = map[string]string{}
    }
    if response.MultiValueHeaders == nil {
        response.MultiValueHeaders = map[string][]string{}
    }
    
    // https://forums.aws.amazon.com/thread.jspa?threadID=94483
    response.IsBase64Encoded = false
    response.StatusDescription = http.StatusText(response.StatusCode)
    response.Headers["Content-Type"] = "application/json; charset=utf-8"
    if response.Body == "" {
        response.Body = "{}"
    }
    
    return response
}

func getLoggingCtx(ctx context.Context, r awsevents.ALBTargetGroupRequest) context.Context {
    
    // If XRAY is enabled, get the trace_id
    traceID := getTraceID(ctx, r)
    
    // Get the customer RequestID or set to UUID if empty
    requestID := getRequestID(r)
    
    // Set to UUID for internal usage
    correlationID := getCorrelationID(r)
    
    rsFields := gcontext.RequestScopedFields{
    	TraceID: traceID, 
    	RequestID: requestID, 
    	CorrelationID: correlationID,
    }
    
    return gcontext.AddRequestFields(ctx, rsFields)
}
```





