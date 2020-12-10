package datadog_test

import (
	"context"
	"github.com/cultureamp/glamplify/log"
	"net/http"
	"testing"

	aws "github.com/aws/aws-lambda-go/events"
	"github.com/cultureamp/glamplify/datadog"
	"gotest.tools/assert"
)

func Test_DataDog_Application(t *testing.T) {
	ctx := context.Background()
	app := datadog.NewApplication(ctx, "Glamplify-Unit-Tests", func(conf *datadog.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = false
	})

	assert.Assert(t, app != nil, "application was nil")

	app.Shutdown()
}

func Test_DataDog_WrapHandler(t *testing.T) {
	ctx := context.Background()
	app := datadog.NewApplication(ctx, "Glamplify-Unit-Tests", func(conf *datadog.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = false
	})

	handler := app.WrapHandler("/", dataDogServeHTTP)
	assert.Assert(t, handler != nil, handler)
}

func dataDogServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func Test_DataDog_WrapLambda(t *testing.T) {
	ctx := context.Background()
	app := datadog.NewApplication(ctx, "Glamplify-Unit-Tests", func(conf *datadog.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = true
	})

	handler := app.WrapLambdaHandler(dataDogLambadaHandler)
	assert.Assert(t, handler != nil, handler)
}

func dataDogLambadaHandler(ctx context.Context, request aws.ALBTargetGroupRequest) (aws.ALBTargetGroupResponse, error) {
	return aws.ALBTargetGroupResponse{StatusCode: 200}, nil
}

func Test_DataDog_RecordMetric(t *testing.T) {
	ctx := context.Background()
	app := datadog.NewApplication(ctx, "Glamplify-Unit-Tests", func(conf *datadog.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = true
	})

	app.RecordLambdaMetric("glamplify-test", 1.0, log.Fields{})
}


func Test_DataDog_TraceHandler(t *testing.T) {
	ctx := context.Background()
	app := datadog.NewApplication(ctx, "Glamplify-Unit-Tests", func(conf *datadog.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = true
	})

	app.WrapLambdaHandler(dataDogLambadaHandler)
	span, _ := app.TraceHandler(ctx, "root")
	assert.Assert(t, span != nil, span)
}