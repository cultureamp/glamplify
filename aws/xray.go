package aws

import (
	"context"
	"github.com/cultureamp/glamplify/env"
	"net/http"

	"github.com/aws/aws-xray-sdk-go/awsplugins/ec2"
	"github.com/aws/aws-xray-sdk-go/awsplugins/ecs"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/aws/aws-xray-sdk-go/xraylog"
)

const (
	// XrayEnv = "XRAY_LOGGING"
	XrayEnv = "XRAY_LOGGING"
)

// TracerConfig for setting initial values for Tracer
type TracerConfig struct {
	Environment   string
	AWSService    string
	EnableLogging bool
	Version       string
}

// Tracer represents an XRAY trace
type Tracer struct {
	config TracerConfig
	logger *xrayLogger
}

// NewTracer creates a new Tracer
func NewTracer(ctx context.Context, configure ...func(*TracerConfig)) *Tracer {
	conf := TracerConfig{
		Environment:   "development",
		EnableLogging: env.GetBool(env.AwsXrayEnv, false),
	}

	for _, config := range configure {
		config(&conf)
	}

	if conf.Environment == "production" {
		if conf.AWSService == "ECS" {
			ecs.Init()
		} else if conf.AWSService == "EC2" {
			ec2.Init()
		}
	}

	logger := newXrayLogger(ctx)
	if conf.EnableLogging {
		xray.SetLogger(logger)
	}

	if err := xray.Configure(xray.Config{ServiceVersion: conf.Version}); err != nil {
		logger.Log(xraylog.LogLevelError, newPrintArgs(err.Error()))
	}

	return &Tracer{
		config: conf,
		logger: logger,
	}
}

// GetTraceID returns the current xray trace_id
func (tracer Tracer) GetTraceID(ctx context.Context) string {
	if xray.RequestWasTraced(ctx) {
		return xray.TraceID(ctx)
	}

	return ""
}

// RoundTripper returns the current xray RoundTripper
func (tracer Tracer) RoundTripper(rt http.RoundTripper) http.RoundTripper {
	return xray.RoundTripper(rt)
}

// SegmentHandler returns the current xray segment http handler
func (tracer Tracer) SegmentHandler(name string, h http.Handler) http.Handler {
	sn := xray.NewFixedSegmentNamer(name)
	return xray.Handler(sn, h)
}

// DynamicSegmentHandler returns the current xray dynamic segment http handler
func (tracer Tracer) DynamicSegmentHandler(fallback string, wildcardHost string, h http.Handler) http.Handler {
	sn := xray.NewDynamicSegmentNamer(fallback, wildcardHost)
	return xray.Handler(sn, h)
}

// Capture wrapper around xray.Capture as per https://docs.aws.amazon.com/xray/latest/devguide/xray-sdk-go-subsegments.html
func (tracer Tracer) Capture(ctx context.Context, name string, fn func(context.Context) error) (err error) {
	return xray.Capture(ctx, name, fn)
}

// AddMetadata wrapper around xray.AddMetadata as per https://docs.aws.amazon.com/xray/latest/devguide/xray-sdk-go-subsegments.html
func (tracer Tracer) AddMetadata(ctx context.Context,  key string, value interface{}) error {
	return xray.AddMetadata(ctx, key, value)
}

// Middleware adds a new XRAY segment when used as a middleware
func (tracer Tracer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sn := xray.NewFixedSegmentNamer(r.URL.Path)
		next = xray.Handler(sn, next)
		next.ServeHTTP(w, r)
	})
}