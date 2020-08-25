package bugsnag

import (
	"context"
	"github.com/bugsnag/bugsnag-go"
	"github.com/cultureamp/glamplify/helper"
	"github.com/cultureamp/glamplify/log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Enabled         bool
	Logging         bool
	License         string
	AppName         string
	AppVersion      string
	ReleaseStage    string
	ProjectPackages []string
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
		Enabled:         false,
		Logging:         false,
		License:         os.Getenv("BUGSNAG_LICENSE_KEY"),
		AppName:         name,
		AppVersion:      helper.GetEnvString("APP_VERSION", "1.0.0"),
		ReleaseStage:    helper.GetEnvString("APP_ENV", "production"),
		ProjectPackages: []string{"github.com/cultureamp"},
	}

	for _, config := range configure {
		config(&conf)
	}

	cfg := bugsnag.Configuration{
		APIKey:          conf.License,
		AppType:         conf.AppName,
		AppVersion:      conf.AppVersion,
		ReleaseStage:    conf.ReleaseStage,
		ProjectPackages: conf.ProjectPackages,
		ParamsFilters:   []string{"password", "pwd"}, // todo - add others
	}

	if conf.Logging {
		cfg.Logger = newBugsnagLogger(context.Background())
	}

	bugsnag.Configure(cfg)

	return &Application{conf: conf}, nil
}

// Shutdown flushes any remaining data to the SAAS endpoint
func (app Application) Shutdown() {
	time.Sleep(waitFORBugsnag)
}

// Adds a Bugsnag when used as middleware
func (app *Application) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = app.addToHTTPContext(r)
		next.ServeHTTP(w, r)
	})
}

func (app *Application) WrapHTTPHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
	p, h := app.wrapHTTPHandler(pattern, http.HandlerFunc(handler))
	return p, func(w http.ResponseWriter, r *http.Request) {
		r = app.addToHTTPContext(r)
		h.ServeHTTP(w, r)
	}
}

func (app *Application) wrapHTTPHandler(pattern string, handler http.Handler) (string, http.Handler) {
	return pattern, bugsnag.Handler(handler)
}

func Error(err error, fields log.Fields) error {
	return internal.Error(err, fields)
}

func (app Application) Error(err error, fields log.Fields) error {
	if !app.conf.Enabled {
		return nil
	}

	ctx := bugsnag.StartSession(context.Background())
	defer bugsnag.AutoNotify(ctx)

	return app.ErrorWithContext(ctx, err, fields)
}

func ErrorWithContext(ctx context.Context, err error, fields log.Fields) error {
	return internal.ErrorWithContext(ctx, err, fields)
}

func (app Application) ErrorWithContext(ctx context.Context, err error, fields log.Fields) error {
	if !app.conf.Enabled {
		return nil
	}

	meta := fieldsAsMetaData(fields)
	return bugsnag.Notify(err, ctx, meta)
}

func (app *Application) addToHTTPContext(req *http.Request) *http.Request {
	ctx := app.addToContext(req.Context())
	return req.WithContext(ctx)
}

func (app *Application) addToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, notifyContextKey, app)
}

func fieldsAsMetaData(fields log.Fields) bugsnag.MetaData {
	meta := make(bugsnag.MetaData)
	for k, v := range fields {
		meta.Add("app context", k, v)
	}
	return meta
}
