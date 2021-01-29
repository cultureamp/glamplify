package datadog

import (
	"context"
	"github.com/cultureamp/glamplify/env"
	"net/http"
	"os"

	ddlambda "github.com/DataDog/datadog-lambda-go"
	gcontext "github.com/cultureamp/glamplify/context"
	"github.com/cultureamp/glamplify/log"
	ddhttp "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Tags are key value pairs used to roll up applications into specific categories
type Tags map[string]string

// Config contains Application and Transaction behavior settings.
// Use NewConfig to create a config with proper defaults.
type Config struct {

	// Enabled controls whether the agent will communicate with the New Relic
	// servers and spawn goroutines.  Setting this to be false is useful in
	// testing and staging situations.
	Enabled bool
	Name    string

	Logging            bool
	APIKey             string
	AppName            string
	AppEnv             string
	AppVersion         string
	AgentHost          string
	AgentStatsDPort    string
	MetricSite         string
	WithAnalytics      bool
	WithRuntimeMetrics bool

	// Tags are global tags applied to all messages sent to Data Dog
	Tags Tags

	// ServerlessMode contains types which control behavior when running in AWS Lambda.
	ServerlessMode bool

	// PRIVATE
	logger *agentLogger
}

// Application is a wrapper over the underlying implementation
type Application struct {
	conf Config
}

// NewApplication creates a new Application - you should only create 1 Application per process
func NewApplication(ctx context.Context, name string, configure ...func(*Config)) *Application {
	// https://docs.datadoghq.com/tracing/setup/go/
	// https://docs.datadoghq.com/getting_started/tagging/unified_service_tagging/?tab=kubernetes
	// We highly recommend using DatadogEnv, DatadogService, and DatadogVersion to set env, service, and version for your services.

	conf := Config{
		Enabled:            false,
		Name:               name,
		Logging:            false,
		APIKey:             os.Getenv(env.DatadogAPIKey),
		AppName:            env.GetString(env.DatadogService, os.Getenv(env.AppNameEnv)),
		AppEnv:             env.GetString(env.DatadogEnv, os.Getenv(env.AppFarmEnv)),
		AppVersion:         env.GetString(env.DatadogVersion, os.Getenv(env.AppVerEnv)),
		AgentHost:          env.GetString(env.DatadogAgentHost, "localhost"),
		AgentStatsDPort:    env.GetString(env.DatadogStatsdPort, "8125"),
		MetricSite:         env.GetString(env.DatadogSite, "datadoghq.com"),
		WithAnalytics:      false,
		WithRuntimeMetrics: false,
		logger:             nil,
	}

	for _, config := range configure {
		config(&conf)
	}

	app := &Application{
		conf: conf,
	}

	if !conf.Enabled {
		return app
	}

	if !conf.ServerlessMode {
		options := []ddtracer.StartOption{
			ddtracer.WithEnv(conf.AppEnv),
			ddtracer.WithService(conf.AppName),
			ddtracer.WithServiceVersion(conf.AppVersion),
			ddtracer.WithAnalytics(conf.WithAnalytics),
			ddtracer.WithDogstatsdAddress(conf.AgentHost + ":" + conf.AgentStatsDPort),
		}

		if conf.WithRuntimeMetrics {
			options = append(options, ddtracer.WithRuntimeMetrics())
		}

		if conf.Logging {
			conf.logger = newAgentLogger(ctx)
			options = append(options, ddtracer.WithLogger(conf.logger))
		}

		if len(conf.Tags) > 0 {
			for k, v := range conf.Tags {
				options = append(options, ddtracer.WithGlobalTag(k, v))
			}
		}

		ddtracer.Start(options...)
	}

	return app
}

// WrapHandler wraps an http.Handler with tracing using the given service and resource.
func (app Application) WrapHandler(resource string, handler http.HandlerFunc) http.Handler {
	if !app.conf.Enabled {
		return handler
	}

	var options []ddhttp.Option
	return ddhttp.WrapHandler(handler, app.conf.AppName, resource, options...)
}

// Shutdown flushes any remaining data to the SAAS endpoint
func (app Application) Shutdown() {
	if app.conf.Enabled && !app.conf.ServerlessMode {
		ddtracer.Stop()
	}
}

// WrapLambdaHandler is used to instrument your lambda functions.
// It returns a modified handler that can be passed directly to the lambda. Start function.
func (app Application) WrapLambdaHandler(handler interface{}) interface{} {
	if app.conf.Enabled {
		c := &ddlambda.Config{
			APIKey:                app.conf.APIKey,
			DebugLogging:          app.conf.Logging,
			Site:                  app.conf.MetricSite,
			EnhancedMetrics:       true,
			ShouldRetryOnFailure:  false,
			ShouldUseLogForwarder: true,
		}

		return ddlambda.WrapHandler(handler, c)
	}

	return handler
}

// RecordLambdaMetric sends a distribution metric to DataDog
func (app Application) RecordLambdaMetric(metricName string, metricValue float64, fields log.Fields) {
	if app.conf.Enabled {
		tags := fields.ToTags(true)
		ddlambda.Metric(
			metricName,
			metricValue,
			tags...,
		)
	}
}

// TraceHandler adds default tags to the current DataDog span and then creates a new child span using the operationName
// This should be called at the top of every AWS Lambda handler(...)
func (app Application) TraceHandler(ctx context.Context, operationName string) (ddtracer.Span, context.Context) {
	span, ok := ddtracer.SpanFromContext(ctx)
	if ok {
		rsFields, ok := gcontext.GetRequestScopedFields(ctx)
		if ok {
			span.SetTag("app_xray_id", rsFields.TraceID)
			span.SetTag("app_request_id", rsFields.RequestID)
			span.SetTag("app_correlation_id", rsFields.CorrelationID)
		}
	}

	return ddtracer.StartSpanFromContext(ctx, operationName)
}

