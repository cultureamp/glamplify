# glamplify
Go Amplify Module of useful common tools. The guiding principle is to implement a very light weight wrapper over the standard library (or where not adequate an open source community library), that conforms to our standard practises (12-Factor) and sensible defaults.


## Install

```
go get github.com/cultureamp/glamplify
```

## Usage

### Config
```Go
package main

import (
	"github.com/cultureamp/glamplify/config"
)

func main() {

    // Settings will contain configuration data as read in from the config file.
    settings := config.Load()

    // Or if you want to look for a config file from a specific location use
    // settings = config.LoadFrom([]string{"${HOME}/settings"}, "config")

    // Then you can use
    if settings.App.Version > 2.0 {
        // to do
    }
}
```
If no config.yml or config.json can be found, or if it is corrupted, then a config will be created by checking these ENV variables.

- CONFIG_APPNAME (default: "service-name")
- CONFIG_VERSION (default: 1.0)

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
		conf.AWSService = "ECS"         // or "EC2" or "LAMBDA"
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
    "github.com/cultureamp/glamplify/aws"
    "github.com/cultureamp/glamplify/constants"
    gcontext "github.com/cultureamp/glamplify/context"
    "github.com/cultureamp/glamplify/jwt"
    "github.com/cultureamp/glamplify/log"
    "net/http"
    "time"
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

### Monitor

Make sure you have the environment variable NEW_RELIC_LICENSE_KEY set to the correct 40 character license key.
Alternatively, you can read it from another environment variable and pass it into the monitor.Config struct.

#### Adding Attributes to a Web Request Transaction
```Go
package main

import (
    "github.com/aws/aws-xray-sdk-go/xray"
    "github.com/cultureamp/glamplify/aws"
    gcontext "github.com/cultureamp/glamplify/context"
    "github.com/cultureamp/glamplify/jwt"
    "github.com/cultureamp/glamplify/log"
    "github.com/cultureamp/glamplify/newrelic"
    "net/http"
)

func main() {

    app, err := monitor.NewApplication("GlamplifyDemo", func(conf *monitor.Config) {
        conf.Enabled = true             // default = "false"
        conf.Logging = true             // default = "false"
        conf.ServerlessMode = false     // default = "false"
     })

    _, handler := app.WrapHTTPHandler("/", requestHandler)
    h := http.HandlerFunc(handler)

    if err = http.ListenAndServe(":8080", h); err != nil {
        panic(err)
    }

    app.Shutdown()
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
    r = gcontext.WrapRequest(r) // Reads and Sets TraceID, RequestID, CorrelationID, User and Customer
    logger := log.NewFromRequest(r)

    // Do things

    txn, err := monitor.TxnFromRequest(w, r)
    if err != nil {
        logger.Error("monitor_transaction_failed", err)
        txn.AddAttributes(log.Fields{
            "aString": "hello world",
            "aInt":    123,
        })
    }

    // Do more things

    if err != nil {
        logger.Error("program_error", err)
        txn.AddAttributes(log.Fields{
            "aString2": "goodbye",
            "aInt2":    456,
        })
    }
}
```

#### Custom Events to a Web Request Transaction
```Go
package main

import (
     "net/http"

    "github.com/aws/aws-xray-sdk-go/xray"
    "github.com/cultureamp/glamplify/aws"
    gcontext "github.com/cultureamp/glamplify/context"
    "github.com/cultureamp/glamplify/jwt"
    "github.com/cultureamp/glamplify/log"
    "github.com/cultureamp/glamplify/newrelic"
)

func main() {

    app, err := monitor.NewApplication("GlamplifyDemo", func(conf *monitor.Config) {
        conf.Enabled = true             // default = "false"
        conf.Logging = true             // default = "false"
        conf.ServerlessMode = false     // default = "false"
    })

    _, handler := app.WrapHTTPHandler("/", requestHandler)
    h := http.HandlerFunc(handler)

    if err = http.ListenAndServe(":8080", h); err != nil {
        panic(err)
    }

    app.Shutdown()
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
    r = gcontext.WrapRequest(r) // Reads and Sets TraceID, RequestID, CorrelationID, User and Customer
    logger := log.NewFromRequest(r)
 
    // Do things
    app, err := monitor.AppFromRequest(w, r)
    if err != nil {
        logger.Error("monitor_failed", err)
    }

    err = app.RecordEvent("mycustomEvent", log.Fields{
        "aString": "hello world",
        "aInt":    123,
    })
    if err != nil {
        logger.Error("monitor_record_failed", err)
    }
    // Do more things
}
```

#### Adding Attributes to a Lambda (Serverless)
```Go
package main

import (
    "context"
    "github.com/aws/aws-xray-sdk-go/xray"
    "github.com/cultureamp/glamplify/aws"
    gcontext "github.com/cultureamp/glamplify/context"
    "github.com/cultureamp/glamplify/jwt"
    "github.com/cultureamp/glamplify/log"
    "github.com/cultureamp/glamplify/newrelic"
    "net/http"
)

func main() {

    app, _ := monitor.NewApplication("GlamplifyDemo", func(conf *monitor.Config) {
        conf.Enabled = true             // default = "false"
        conf.Logging = true             // default = "false"
        conf.ServerlessMode = true      // default = "false"
     })

    monitor.Start(handler, app)
}

func handler(ctx context.Context) {
    ctx = gcontext.WrapCtx(ctx) // reads and sets TraceID and CorrelationID
    logger := log.NewFromCtx(ctx)

    // Do things

    txn, err := monitor.TxnFromContext(ctx)
    if err != nil {
        logger.Error("monitor_transaction_failed", err)
        txn.AddAttributes(log.Fields{
            "aString": "hello world",
            "aInt":    123,
        })
    }

    // Do more things

    if err != nil {
        logger.Error("program_failed", err)
        txn.AddAttributes(log.Fields{
            "aString2": "goodbye",
            "aInt2":    456,
        })
    }
}
```

#### Custom Events to a Lambda (Serverless)
```Go
package main

import (
    "context"
    "github.com/aws/aws-xray-sdk-go/xray"
    "github.com/cultureamp/glamplify/aws"
    gcontext "github.com/cultureamp/glamplify/context"
    "github.com/cultureamp/glamplify/jwt"
    "github.com/cultureamp/glamplify/log"
    "github.com/cultureamp/glamplify/newrelic"
    "net/http"
)

func main() {

    app, _ := monitor.NewApplication("GlamplifyDemo", func(conf *monitor.Config) {
        conf.Enabled = true             // default = "false"
        conf.Logging = true             // default = "false"
        conf.ServerlessMode = true      // default = "false"
    })

    monitor.Start(handler, app)
}

func handler(ctx context.Context) {
    ctx = gcontext.WrapCtx(ctx)  // reads and sets TraceID and CorrelationID
    logger := log.NewFromCtx(ctx)

    // Do things

    app, err := monitor.AppFromContext(ctx)
    if err != nil {
        logger.Error("monitor_transaction_failed", err)
        err = app.RecordEvent("mycustomEvent", log.Fields{
            "aString": "hello world",
            "aInt":    123,
        })
    }
    // Do more things
}
```

## Notify

Make sure you have the environment variable BUGSNAG_LICENSE_KEY set to the correct license key.
Alternatively, you can read it from another environment variable and pass it into the notify.Config struct.

```Go
package main

import (
    "errors"
    gcontext "github.com/cultureamp/glamplify/context"
    "github.com/cultureamp/glamplify/jwt"
    "github.com/cultureamp/glamplify/log"
    "github.com/cultureamp/glamplify/bugsnag"
    "net/http"
)

func main() {

    notifier, err := notify.NewNotifier("GlamplifyDemo", func(conf *notify.Config) {
        conf.Enabled = true             // default = "false"
        conf.Logging = true             // default = "false"
     })

    _, handler := notifier.WrapHTTPHandler("/", requestHandler)
    h := http.HandlerFunc(handler)

    if err = http.ListenAndServe(":8080", h); err != nil {
        panic(err)
    }

    notifier.Shutdown()
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
    r = gcontext.WrapRequest(r)  // Reads and Sets TraceID, RequestID, CorrelationID, User and Customer
    logger := log.NewFromRequest(r)

    notifier, notifyErr := notify.NotifyFromRequest(w, r)
    if notifyErr != nil {
        logger.Error("notification_error", notifyErr)
    }

    // Do things

    // Pretend we got an error
    err := errors.New("NPE")
    notifier.ErrorWithContext(r.Context(), err, log.Fields {
        "key": "value",
    })
}
```

