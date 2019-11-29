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

    // settings will contain configuration data as read in from the config file.
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

### Logging

```Go
package main

import (
    "bytes"
    "errors"

    "github.com/cultureamp/glamplify/log"
)

func main() {

    // You can either get a new logger, or just use the public functions which internally use an internal logger
    // eg. log.Debug(), log.Print() and log.Error()

    // Example below shows usage with the package level logger (sensible default), but can 
    // use an instance of a logger by calling mylogger := log.New()

    // Emit debug trace
    // All messages must be static strings (as per Culture Amp Sensibile Default)
    log.Debug("Something happened")

    field
    // Fields can contain any type of variables
    log.Debug("Something happened", log.Fields{
        "aString": "hello",
        "aInt":    123,
        "aFloat":  42.48,
     })

    field
    // Typically Print will be sent onto 3rd party aggregation tools (eg. Splunk)
    log.Print("Executing main")

    field
    // Fields can contain any type of variables
    log.Print("Executing main", log.Fields{
        "program-name": "helloworld.exe",
        "start-up-param":    123,
        "user":  "admin",
    })

    field
    // Errors will always be sent onto 3rd party aggregation tools (eg. Splunk)
    err := errors.New("missing database connection string")
    log.Error(err, "Main program stopped unexpectedly")

    field
    // Fields can contain any type of variables
    err := errors.New("missing database connection string")
    log.Error(err, "Executing main", log.Fields{
        "program-name": "helloworld.exe",
        "start-up-param":    123,
        "user":  "admin",
     })

    field
    // have a requestID for every log message within that scope) then you can use WithScope()
    scope := log.WithScope(log.Fields { "requestID" : 123 })

    // then just use the scope as you would a normal logger
    field
    field
    scope.Print("Starting web request", log.Fields { "auth": "oauth" })

    // If you want to change the output or time format you can only do this for an
    // instance of the logger you create (not the internal one) by doing this:

    memBuffer := &bytes.Buffer{}
    logger := log.New(func(conf *log.Config) {
        conf.Output = memBuffer                     // can be set to anything that support io.Write
        conf.TimeFormat = "2006-01-02T15:04:05"     // any valid time format
        conf.debugForwardLogTo = "splunk"           // send debug messages to "splunk"
    })

    // The internall logger will always use these default values:
    // output = os.Stdout
    // time format = "2006-01-02T15:04:05.000Z07:00"
    // debugForwardLogTo = "none"
    // printForwardLogTo = "splunk"
    // errorForwardLogTo = "splunk"
}

```
Use `log.Debug` for logging that will only be used when diving deep to uncover bugs. Typically `log.Debug` messages will not automatically be sent to other 3rd party systems (eg. Splunk).

Use `log.Print` for standard log messages that you want to see always. These will never be turned off and will likely be always sent to 3rd party systems for further analysis (eg. Spliunk).

Use `log.Error` when you have encounter a GO error. This will NOT stop the program, it is up to you to call exit() or panic() if this is not recoverable. All error messages will be forwarded to 3rd party systems for monitoring and further analysis.

### Monitor

Make sure you have the environment variable NEW_RELIC_LICENSE_KEY set to the correct 40 character license key. Alternatively, you can read it from another environment variable and pass it into the event.Config struct.

#### Adding Attributes to a Web Request Transaction
```Go
package main

import (
    "net/http"
    "github.com/cultureamp/glamplify/monitor"
)

func main() {

    app, err := monitor.NewApplication("GlamplifyDemo", func(conf *monitor.Config) {
        conf.Enabled = true             // default = "false"
        conf.Logging = true             // default = "false"
        conf.ServerlessMode = false     // default = "false"
     })

    _, handler := app.WrapTxnHandler("/", rootRequestHandler)
    http.HandlerFunc(handler)

    if err := http.ListenAndServe(":8080", nil); err != nil {
        panic(err)
    }

    app.Shutdown()
}

func rootRequestHandler(w http.ResponseWriter, r *http.Request) {

    // Do things

    txn, ok := monitor.TxnFromRequest(w, r)
    if ok {
        txn.AddAttributes(monitor.Fields{
            "aString": "hello world",
            "aInt":    123,
        })
    }

    // Do more things

    if ok {
        txn.AddAttributes(monitor.Fields{
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
    "github.com/cultureamp/glamplify/monitor"
)

func main() {

    app, err := monitor.NewApplication("GlamplifyDemo", func(conf *monitor.Config) {
		conf.Enabled = true             // default = "false"
		conf.Logging = true             // default = "false"
		conf.ServerlessMode = false     // default = "false"
	})

    _, handler := app.WrapTxnHandler("/", rootRequestHandler)
    http.HandlerFunc(handler)

    if err := http.ListenAndServe(":8080", nil); err != nil {
        panic(err)
    }
    
    app.Shutdown()
}

func rootRequestHandler(w http.ResponseWriter, r *http.Request) {

    // Do things

    err = monitor.RecordEvent("mycustomEvent", monitor.Fields{
        "aString": "hello world",
        "aInt":    123,
    })

    // Do more things
}
```

#### Adding Attributes to a Lambda (Serverless)
```Go
package main

import (
    "context"
    "net/http"
    "github.com/cultureamp/glamplify/monitor"
)

func main() {
    app, err := monitor.NewApplication("GlamplifyDemo", func(conf *monitor.Config) {
        conf.Enabled = true             // default = "false"
        conf.Logging = true             // default = "false"
        conf.ServerlessMode = true      // default = "false"
     })

    monitor.Start(handler, app)
}

func handler(ctx context.Context) {

    // Do things

    txn, err := monitor.TxnFromContext(ctx)
    if err != nil {
        txn.AddAttributes(monitor.Fields{
            "aString": "hello world",
            "aInt":    123,
        })
    }

    // Do more things

    if err != nil {
        txn.AddAttributes(monitor.Fields{
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
    "net/http"
    "github.com/cultureamp/glamplify/monitor"
)

func main() {
    app, err := monitor.NewApplication("GlamplifyDemo", func(conf *monitor.Config) {
        conf.Enabled = true             // default = "false"
        conf.Logging = true             // default = "false"
        conf.ServerlessMode = true      // default = "false"
    })

    monitor.Start(handler, app)
}

func handler(ctx context.Context) {

    // Do things

    app, err := monitor.AppFromContext(ctx)
    if err != nil {
        err = app.RecordEvent("mycustomEvent", monitor.Fields{
            "aString": "hello world",
            "aInt":    123,
        })

    // Do more things
}
```