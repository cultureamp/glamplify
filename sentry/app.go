package sentry

import (
	"context"
	"github.com/cultureamp/glamplify/env"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

// Config represents Sentry configuration values
type Config struct {
	Enabled          bool
	Logging          bool
	DSN              string
	Transport        sentry.Transport
	FlushTimeoutInMs int
	ServerName       string
	AppName          string
	AppVersion       string
	ReleaseStage     string
}

// Application represents a sentry app
type Application struct {
	conf Config
}

// NewApplication creates a new sentry Application
func NewApplication(ctx context.Context, name string, configure ...func(*Config)) (*Application, error) {

	if len(name) == 0 {
		name = env.GetString(env.AppNameEnv, "default")
	}

	conf := Config{
		Enabled:          false,
		Logging:          false,
		DSN:              os.Getenv(env.SentryDsnEnv),
		FlushTimeoutInMs: env.GetInt(env.SentryFlushTimeoutInMsEnv, 500),
		Transport:        sentry.NewHTTPTransport(),
		AppName:          name,
		AppVersion:       env.GetString(env.AppVerEnv, "1.0.0"),
		ReleaseStage:     env.GetString(env.AppFarmEnv, "production"),
	}

	host, err := os.Hostname()
	if err == nil {
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
		Transport:   conf.Transport,
		Release:     conf.AppName + "@" + conf.AppVersion,
		Environment: conf.ReleaseStage,
		ServerName:  conf.ServerName,
	}

	if conf.Logging {
		cfg.Debug = true
		cfg.DebugWriter = newSentryLogger(ctx)
	}

	err = sentry.Init(cfg)
	return &Application{conf: conf}, err
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
// CaptureException or CaptureMessage.
//
// For Serverless Lambda either call this (with a small duration) at the end of
// each request, or configure Sentry to send events over the network synchronously,
// configure it to use the HTTPSyncTransport.
func (app Application) Flush(timeout time.Duration) bool {
	if app.conf.Enabled {
		return sentry.Flush(timeout)
	}

	return true
}

// FlushDefault waits until the underlying Transport sends any buffered events to the
// Sentry server, blocking for at most the given timeout. It returns false if
// the timeout was reached. In that case, some events may not have been sent.
//
// FlushDefault should be called before terminating the program to avoid
// unintentionally dropping events.
//
// Do not call FlushDefault indiscriminately after every call to CaptureEvent,
// CaptureException or CaptureMessage.
//
// For Serverless Lambda either call this (with a small duration) at the end of
// each request, or configure Sentry to send events over the network synchronously,
// configure it to use the HTTPSyncTransport.
func (app Application) FlushDefault() bool {
	if app.conf.Enabled {
		return sentry.Flush(time.Duration(app.conf.FlushTimeoutInMs) * time.Millisecond)
	}

	return true
}

// Middleware adds Sentry as a middleware
func (app *Application) Middleware(next http.Handler) http.Handler {
	if !app.conf.Enabled {
		return next
	}

	// Create an instance of sentryhttp
	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	return sentryHandler.Handle(next)
}

// WrapHTTPHandler wraps a http handler with Sentry
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

// Error records and sends an error with sentry
func (app Application) Error(err error) *sentry.EventID {
	if !app.conf.Enabled {
		return nil
	}

	return sentry.CaptureException(err)
}

// Message records and sends a meesage with sentry
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
