package notify

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
	Enabled         bool     `yaml:"enabled"`
	Logging         bool     `yaml:"logging"`
	License         string   `yaml:"license"`
	AppName         string   `yaml:"app_name"`
	AppVersion      string   `yaml:"app_version"`
	ReleaseStage    string   `yaml:"release_stage"`
	ProjectPackages []string `yaml:"proejct_packages"`
}

type Notifier struct {
	conf Config
}

const (
	waitFORBugsnag = 2 * time.Second
)

var (
	internal, _ = NewNotifier(helper.GetEnvString("APP_NAME", "default"), func(conf *Config) { conf.Enabled = true })
)

func NewNotifier(name string, configure ...func(*Config)) (*Notifier, error) {

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
		cfg.Logger = newNotifyLogger(context.Background())
	}

	bugsnag.Configure(cfg)

	return &Notifier{conf: conf}, nil
}

// Shutdown flushes any remaining data to the SAAS endpoint
func (notify Notifier) Shutdown() {
	time.Sleep(waitFORBugsnag)
}

// Adds a Bugsnag when used as middleware
func (notify *Notifier) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = notify.addToHTTPContext(r)
		next.ServeHTTP(w, r)
	})
}

func (notify *Notifier) WrapHTTPHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
	p, h := notify.wrapHTTPHandler(pattern, http.HandlerFunc(handler))
	return p, func(w http.ResponseWriter, r *http.Request) {
		r = notify.addToHTTPContext(r)
		h.ServeHTTP(w, r)
	}
}

func (notify *Notifier) wrapHTTPHandler(pattern string, handler http.Handler) (string, http.Handler) {
	return pattern, bugsnag.Handler(handler)
}

func Error(err error, fields log.Fields) error {
	return internal.Error(err, fields)
}

func (notify Notifier) Error(err error, fields log.Fields) error {
	if !notify.conf.Enabled {
		return nil
	}

	ctx := bugsnag.StartSession(context.Background())
	defer bugsnag.AutoNotify(ctx)

	return notify.ErrorWithContext(ctx, err, fields)
}

func ErrorWithContext(ctx context.Context, err error, fields log.Fields) error {
	return internal.ErrorWithContext(ctx, err, fields)
}

func (notify Notifier) ErrorWithContext(ctx context.Context, err error, fields log.Fields) error {
	if !notify.conf.Enabled {
		return nil
	}

	meta := fieldsAsMetaData(fields)
	return bugsnag.Notify(err, ctx, meta)
}

func (notify *Notifier) addToHTTPContext(req *http.Request) *http.Request {
	ctx := notify.addToContext(req.Context())
	return req.WithContext(ctx)
}

func (notify *Notifier) addToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, notifyContextKey, notify)
}

func fieldsAsMetaData(fields log.Fields) bugsnag.MetaData {
	meta := make(bugsnag.MetaData)
	for k, v := range fields {
		meta.Add("app context", k, v)
	}
	return meta
}
