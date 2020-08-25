package newrelic

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/cultureamp/glamplify/log"

	newrelic "github.com/newrelic/go-agent"
)

const (
	waitFORNR = 4 * time.Second
)

// Labels are key value pairs used to roll up applications into specific categories
type Labels map[string]string

// config contains Application and Transaction behavior settings.
// Use NewConfig to create a config with proper defaults.
type Config struct {

	// Enabled controls whether the agent will communicate with the New Relic
	// servers and spawn goroutines.  Setting this to be false is useful in
	// testing and staging situations.
	Enabled bool

	// license is your New Relic license key.
	//
	// https://docs.newrelic.com/docs/accounts/install-new-relic/account-setup/license-key
	License string

	// Logging controls whether Event logging is sent to StdOut or not
	Logging bool

	// Labels are key value pairs used to roll up applications into specific categories.
	//
	// https://docs.newrelic.com/docs/using-new-relic/user-interface-functions/organize-your-data/labels-categories-organize-apps-monitors
	Labels Labels

	// ServerlessMode contains types which control behavior when running in
	// AWS Lambda.
	//
	// https://docs.newrelic.com/docs/serverless-function-monitoring/aws-lambda-monitoring/get-started/introduction-new-relic-monitoring-aws-lambda
	ServerlessMode bool

	// coreLogger logger
	logger *monitorLogger
}

// Application is a wrapper over the underlying implementation
type Application struct {
	impl newrelic.Application
	conf Config
}

// NewApplication creates a new Application - you should only create 1 Application per process
func NewApplication(name string, configure ...func(*Config)) (*Application, error) {

	conf := Config{
		Enabled:        false,
		Logging:        false,
		License:        os.Getenv("NEW_RELIC_LICENSE_KEY"),
		ServerlessMode: false,
		logger:         nil,
	}

	for _, config := range configure {
		config(&conf)
	}

	cfg := newrelic.NewConfig(name, conf.License)
	cfg.Enabled = conf.Enabled // useful to turn on/off in test/dev vs production accounts
	cfg.License = conf.License
	cfg.CustomInsightsEvents.Enabled = true // otherwise custom events won't fire
	cfg.ErrorCollector.Enabled = true
	cfg.ErrorCollector.CaptureEvents = true
	cfg.HighSecurity = false // HighSecurity blocks sending custom events
	cfg.Labels = conf.Labels // camp, environment, data classification, etc
	cfg.RuntimeSampler.Enabled = true
	cfg.ServerlessMode.Enabled = conf.ServerlessMode
	cfg.TransactionTracer.Enabled = true
	cfg.Utilization.DetectAWS = true
	cfg.Utilization.DetectDocker = true

	// for now we turn off DistributedTracing because it is too expensive
	cfg.DistributedTracer.Enabled = false

	if conf.Logging {
		//cfg.Logger = newrelic.NewDebugLogger(os.Stdout) <- this writes JSON to Stdout :(
		// So we have our own implementation that wraps our standard logger
		// TODO - the context here is missing traceID, correlationID, etc
		conf.logger = newMonitorLogger(context.Background())
		cfg.Logger = conf.logger

		cfg.Logger.Debug("configuration", log.Fields{
			"enabled":         conf.Enabled,
			"logging":         conf.Logging,
			"labels":          conf.Labels,
			"serverless_mode": conf.ServerlessMode,
		})
	}

	app := &Application{
		conf: conf,
	}

	impl, err := newrelic.NewApplication(cfg)
	if err != nil {
		app.logError("Failed to create Application", err)
		return nil, err
	}

	if !conf.ServerlessMode {
		// if conf.ServerlessMode = false (server mode) then newrelic.NewApplication spins up
		// some go routines that make a network call back to NR. Until this happens any "RecordCustomEvents"
		// seem to get dropped!
		// Waiting here so that everything is set up and ready
		time.Sleep(waitFORNR)
	}

	app.impl = impl
	return app, err
}

// RecordEvent sends a custom event with the associated data to the underlying implementation
func (app Application) RecordEvent(eventType string, fields log.Fields) error {
	app.log("record_event_begin", log.Fields{"event_type": eventType}, fields)

	// NewRelic has limits on number and size of entries
	// https://docs.newrelic.com/docs/insights/insights-data-sources/custom-data/insights-custom-data-requirements-limits
	// However, if you pass in a string entry longer than 255 it fails "siliently"!!!!!
	// TODO - implement our own checking?

	ok, err := fields.ValidateNewRelic()
	if !ok {
		app.logError("record_event_error", err)
		app.log("record_event_end", log.Fields{"event_type": eventType}, fields)
		return err
	}

	err = app.impl.RecordCustomEvent(eventType, fields)
	app.logError("record_event_error", err)
	app.log("record_event_end", log.Fields{"event_type": eventType}, fields)

	return err
}

// Adds a new NR transaction when used as middleware
func (app *Application) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		txn := app.startTransaction(r.URL.Path, w, r)
		defer txn.End()

		r = txn.addToHTTPContext(r)
		next.ServeHTTP(txn, r)
	})
}

// WrapHTTPHandler adds a Transaction within the current request
func (app *Application) WrapHTTPHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
	p, h := app.wrapHTTPHandler(pattern, http.HandlerFunc(handler))
	return p, func(w http.ResponseWriter, r *http.Request) { h.ServeHTTP(w, r) }
}

// Shutdown flushes any remaining data to the SAAS endpoint
func (app Application) Shutdown() {

	if !app.conf.ServerlessMode {
		// if conf.ServerlessMode = false (server mode) then newrelic.Shutdown can exit its coreLogger go routines
		// before it has sent all pending data!
		// Waiting here so that everything is sent before we start closing down...
		time.Sleep(waitFORNR)
	}

	// The time duration passed here is how long to wait before the shutdown channel processes the request
	// It is NOT how long to wait to send data before shutting down.
	app.impl.Shutdown(30 * time.Second)
}

func (app *Application) wrapHTTPHandler(pattern string, handler http.Handler) (string, http.Handler) {
	return pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		txn := app.startTransaction(pattern, w, r)
		defer txn.End()

		r = txn.addToHTTPContext(r)
		handler.ServeHTTP(txn, r)
	})
}

func (app *Application) startTransaction(name string, w http.ResponseWriter, r *http.Request) Transaction {
	app.log("transaction_start", log.Fields{
		"txn_name": name,
	})

	// Create our wrapper txn and add it to the HTTP ctx
	txn := Transaction{
		app:     app,
		name:    name,
		logging: app.conf.Logging,
		logger:  app.conf.logger,
	}

	// call the NR implementation
	impl := app.impl.StartTransaction(name, w, r)
	txn.impl = impl

	return txn
}

func (app *Application) addToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, appContextKey, app)
}

func (app Application) log(msg string, logFields log.Fields, fields ...log.Fields) {
	if app.conf.Logging {
		merged := logFields.Merge(fields...)
		app.conf.logger.Debug(msg, merged)
	}
}

func (app Application) logError(msg string, err error) {
	if err != nil && app.conf.Logging {
		app.conf.logger.Error(msg, map[string]interface{}{
			"error": err,
		})
	}
}
