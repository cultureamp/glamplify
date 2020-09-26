package datadog

import (
	"context"
	ddlambda "github.com/DataDog/datadog-lambda-go"
	"github.com/cultureamp/glamplify/helper"
	"github.com/cultureamp/glamplify/log"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"os"
)

// Tags are key value pairs used to roll up applications into specific categories
type Tags map[string]string

// config contains Application and Transaction behavior settings.
// Use NewConfig to create a config with proper defaults.
type Config struct {

	// Enabled controls whether the agent will communicate with the New Relic
	// servers and spawn goroutines.  Setting this to be false is useful in
	// testing and staging situations.
	Enabled bool
	Name    string

	Logging            bool
	ApiKey             string
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
	// We highly recommend using DDEnv, DDService, and DDVersion to set env, service, and version for your services.

	conf := Config{
		Enabled:            false,
		Name:               name,
		Logging:            false,
		ApiKey:             os.Getenv(DDApiKey),
		AppName:            helper.GetEnvString(DDService, os.Getenv(log.AppNameEnv)),
		AppEnv:             helper.GetEnvString(DDEnv, os.Getenv(log.AppFarmEnv)),
		AppVersion:         helper.GetEnvString(DDVersion, os.Getenv(log.AppVerEnv)),
		AgentHost:          helper.GetEnvString(DDAgentHost, "localhost"),
		AgentStatsDPort:    helper.GetEnvString(DDDogStatsdPort, "8125"),
		MetricSite:         helper.GetEnvString(DDSite, "datadoghq.com"),
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

// Shutdown flushes any remaining data to the SAAS endpoint
func (app Application) Shutdown() {

	if !app.conf.ServerlessMode {
		ddtracer.Stop()
	}
}

func (app Application) WrapLambdaHandler(handler interface{}) interface{} {
	c := &ddlambda.Config{
		APIKey:       app.conf.ApiKey,
		DebugLogging: app.conf.Logging,
		Site:         app.conf.MetricSite,

		// TODO support other config values?
		ShouldRetryOnFailure:  false,
		ShouldUseLogForwarder: true,
	}

	return ddlambda.WrapHandler(handler, c)
}

func (app Application) RecordLambdaMetric(metricName string, metricValue float64, fields log.Fields) {
	tags := fields.ToTags(true)
	ddlambda.Metric(
		metricName,
		metricValue,
		tags...,
	)
}
