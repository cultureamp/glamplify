package sentry

import (
	"context"
	"github.com/cultureamp/glamplify/log"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"net/http"
	"os"
	"time"

	"github.com/cultureamp/glamplify/helper"
	"github.com/getsentry/sentry-go"
)

type Config struct {
	Enabled      bool
	Logging      bool
	DSN          string
	ServerName   string
	AppName      string
	AppVersion   string
	ReleaseStage string
}

type Application struct {
	conf Config
}

func NewApplication(ctx context.Context, name string, configure ...func(*Config)) (*Application, error) {

	if len(name) == 0 {
		name = helper.GetEnvString(log.AppNameEnv, "default")
	}

	conf := Config{
		Enabled:      false,
		Logging:      false,
		DSN:          os.Getenv("SENTRY_DSN"),
		AppName:      name,
		AppVersion:   helper.GetEnvString(log.AppVerEnv, "1.0.0"),
		ReleaseStage: helper.GetEnvString(log.AppFarmEnv, "production"),
	}

	host, err := os.Hostname()
	if err != nil {
		conf.ServerName = host
	}

	for _, config := range configure {
		config(&conf)
	}

	if !conf.Enabled {
		// if not enabled, then return early...
		return &Application{conf: conf}, nil
	}

	cfg := sentry.ClientOptions{
		Dsn:         conf.DSN,
		Release:     conf.AppName + "@" + conf.AppVersion,
		Environment: conf.ReleaseStage,
		ServerName:  conf.ServerName,
	}

	if conf.Logging {
		cfg.Debug = true
		cfg.DebugWriter = newSentryLogger(ctx)
	}

	err = sentry.Init(cfg)
	if err != nil {
		return nil, err
	}

	return &Application{conf: conf}, nil
}

// Shutdown flushes any remaining data to the SAAS endpoint
func (app Application) Shutdown() {
	if app.conf.Enabled {
		sentry.Flush(2 * time.Second)
	}
}

// Flush waits until the underlying Transport sends any buffered events to the
// Sentry server, blocking for at most the given timeout. It returns false if
// the timeout was reached. In that case, some events may not have been sent.
//
// Flush should be called before terminating the program to avoid
// unintentionally dropping events.
//
// Do not call Flush indiscriminately after every call to CaptureEvent,
// CaptureException or CaptureMessage. Instead, to have the SDK send events over
// the network synchronously, configure it to use the HTTPSyncTransport in the
// call to Init.
func (app Application) Flush(timeout time.Duration) {
	if app.conf.Enabled {
		sentry.Flush(timeout)
	}
}

// Adds Sentry as a middleware
func (app *Application) Middleware(next http.Handler) http.Handler {
	if !app.conf.Enabled {
		return next
	}

	// Create an instance of sentryhttp
	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	return sentryHandler.Handle(next)
}

func (app *Application) WrapHTTPHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
	if !app.conf.Enabled {
		return pattern, handler
	}

	p, h := app.wrapHTTPHandler(pattern, http.HandlerFunc(handler))
	return p, func(w http.ResponseWriter, r *http.Request) {
		r = app.addToHTTPContext(r)
		h.ServeHTTP(w, r)
	}
}

func (app *Application) wrapHTTPHandler(pattern string, handler http.Handler) (string, http.Handler) {
	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	return pattern, sentryHandler.Handle(handler)
}

func (app Application) Error(err error) *sentry.EventID {
	if !app.conf.Enabled {
		return nil
	}

	return sentry.CaptureException(err)
}

func (app Application) Message(message string) *sentry.EventID {
	if !app.conf.Enabled {
		return nil
	}

	return sentry.CaptureMessage(message)
}

func (app *Application) addToHTTPContext(req *http.Request) *http.Request {
	ctx := app.addToContext(req.Context())
	return req.WithContext(ctx)
}

func (app *Application) addToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, sentryContextKey, app)
}
