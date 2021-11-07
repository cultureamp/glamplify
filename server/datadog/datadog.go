package datadog

import (
	"context"
	"github.com/cultureamp/glamplify/settings"
	"glamplify/server/auth"
	"net/http"
	"strconv"
	"time"

	gcontext "github.com/cultureamp/glamplify/context"
	"github.com/urfave/negroni"

	ddlambda "github.com/DataDog/datadog-lambda-go"
	"github.com/cultureamp/glamplify/aws"
	ddclient "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	configured bool
	tags       [][]string
)

// Configure Datadog for the Lambda
func Configure(ctx context.Context, settings settings.Settings) {
	if !settings.DatadogEnabled {
		return
	}

	// Xray is considered part of the Datadog configuration, as
	// the local xray active tracing service is used by datadog to add information,
	// even if datadog isn't configured to pull Xray traces out of AWS. Also
	// we're not typically using Xray without DD.
	_ = aws.NewTracer(ctx, func(config *aws.TracerConfig) {
		config.EnableLogging = false
		config.Version = settings.AppVersion
		config.Environment = settings.AppEnv
		config.AWSService = "lambda"
	})

	// tags to be added to each trace
	tags = [][]string{
		{"farm", settings.Farm},
	}

	configured = true
}

// WrapHandler enables Datadog tracing on the Lambda handler
func WrapHandler(handler interface{}) interface{} {
	c := &ddlambda.Config{
		DebugLogging:          false,
		EnhancedMetrics:       true,
		ShouldUseLogForwarder: true,
	}

	return ddlambda.WrapHandler(handler, c)
}

func RecordRequestStart(ctx context.Context, httpMethod string, url string) {
	span, ok := ddtracer.SpanFromContext(ctx)
	if ok {
		span.SetTag("http.method", httpMethod)
		span.SetTag("http.url", url)

		// ensure cold start tag is present in the request trace as well
		coldStart := ctx.Value("cold_start")
		if coldStart != nil {
			span.SetTag("cold_start", coldStart)
		}

		// add default tags
		for _, tag := range tags {
			span.SetTag(tag[0], tag[1])
		}
	}
}

func RecordRequestComplete(ctx context.Context, httpStatus int) {
	span, ok := ddtracer.SpanFromContext(ctx)
	if ok {
		span.SetTag("http.status_code", strconv.Itoa(httpStatus))
	}
}

// Middleware creates HTTP middleware suitable for tracing calls to the Lambda
func Middleware(next http.Handler, serviceName string) http.Handler {
	if !configured {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		span, ctx := startHandlerTrace(req.Context(), serviceName, req.Method+" "+req.URL.Path)
		defer span.Finish()

		RecordRequestStart(ctx, req.Method, req.URL.Path)

		wrappedWriter := negroni.NewResponseWriter(w)
		next.ServeHTTP(wrappedWriter, req.WithContext(ctx))

		RecordRequestComplete(ctx, wrappedWriter.Status())
	})
}

func startHandlerTrace(ctx context.Context, serviceName string, operationName string) (ddtracer.Span, context.Context) {

	// start a span representing the HTTP request being made
	span, ctx := ddtracer.StartSpanFromContext(ctx, operationName, ddtracer.ServiceName(serviceName))

	// add cross-service tracing fields for context
	if rsFields, ok := gcontext.GetRequestScopedFields(ctx); ok {
		span.SetTag("app.request_id", rsFields.RequestID)
		span.SetTag("app.correlation_id", rsFields.CorrelationID)

		if jwt, ok := auth.GetJWTPayload(ctx); ok {
			span.SetTag("id.user", jwt.Payload.EffectiveUser)
			span.SetTag("id.user-real", jwt.Payload.RealUser)
			span.SetTag("id.customer", jwt.Payload.Customer)
		}
	}

	return span, ctx
}

// WrapHTTPClient wraps a request with Datadog tracing
func WrapHTTPClient(wrapped *http.Client) *http.Client {
	// first wrap with an http client trace metric recorder
	wrapped = withClientTraceRecorder(wrapped)

	// wrap the lot with the datadog handler
	return ddclient.WrapClient(wrapped)
}

// Trace adds a trace of the given name to the current span
// and returns a function to be called when the operation is complete.
func Trace(ctx context.Context, operationName string) (spanCtx context.Context, finishFunc func()) {
	span, spanCtx := ddtracer.StartSpanFromContext(ctx, operationName)
	finishFunc = func() { span.Finish() }

	return
}

// Tag will add data to the current span using
// the supplied name and value.
func Tag(ctx context.Context, tagName string, tagValue interface{}) {
	if span, ok := ddtracer.SpanFromContext(ctx); ok {
		span.SetTag(tagName, tagValue)
	}
}

// TagExecutionTiming will add function timing to the current span using
// the supplied name. The returned function causes the time to recorded and
// is expected to be passed to `defer`.
func TagExecutionTiming(ctx context.Context, tagName string) func() {
	start := time.Now()

	return func() {
		duration := time.Since(start)

		if span, ok := ddtracer.SpanFromContext(ctx); ok {
			span.SetTag(tagName, duration.Milliseconds())
		}
	}
}
