package bugsnag

import (
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
	License      string
	AppName      string
	AppVersion   string
	ReleaseStage string
}

type Application struct {
	conf Config
}

const (
	waitFORBugsnag = 2 * time.Second
)

var (
	internal, _ = NewApplication(helper.GetEnvString("APP_NAME", "default"), func(conf *Config) { conf.Enabled = true })
)

func NewApplication(name string, configure ...func(*Config)) (*Application, error) {

	if len(name) == 0 {
		name = helper.GetEnvString("APP_NAME", "default")
	}

	conf := Config{
		Enabled:      false,
		Logging:      false,
		License:      os.Getenv("SENTRY_LICENSE_KEY"),
		AppName:      name,
		AppVersion:   helper.GetEnvString("APP_VERSION", "1.0.0"),
		ReleaseStage: helper.GetEnvString("APP_ENV", "production"),
	}

	for _, config := range configure {
		config(&conf)
	}

	if conf.Enabled {
		// todo
		// set DSN
	}

	cfg := sentry.ClientOptions{
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		// Debug: conf.Logging,

		// TODO
	}

	err := sentry.Init(cfg)
	if err != nil {
		return nil, err
	}

	return &Application{conf: conf}, nil
}

// Shutdown flushes any remaining data to the SAAS endpoint
func (app Application) Shutdown() {
	sentry.Flush(2 * time.Second)
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
	sentry.Flush(timeout)
}

// Adds a Bugsnag when used as middleware
func (app *Application) Middleware(next http.Handler) http.Handler {

	// Create an instance of sentryhttp
	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	return sentryHandler.Handle(next)
}

func (app *Application) WrapHTTPHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	return pattern, sentryHandler.HandleFunc(handler)
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

