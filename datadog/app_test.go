package datadog_test

import (
	"context"
	"github.com/cultureamp/glamplify/log"
	"net/http"
	"testing"

	aws "github.com/aws/aws-lambda-go/events"
	"github.com/cultureamp/glamplify/datadog"
	"github.com/stretchr/testify/assert"
)

func Test_DataDog_NewApplication(t *testing.T) {
	ctx := context.Background()

	// not enabled
	app := datadog.NewApplication(ctx, "Glamplify-Unit-Tests", func(conf *datadog.Config) {
		conf.Enabled = false
		conf.Logging = false
		conf.ServerlessMode = false
	})
	assert.NotNil(t, app)
	app.Shutdown()

	// serverless
	app = datadog.NewApplication(ctx, "Glamplify-Unit-Tests", func(conf *datadog.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = true
	})
	assert.NotNil(t, app)
	app.Shutdown()

	// EC2/Fargate
	app = datadog.NewApplication(ctx, "Glamplify-Unit-Tests", func(conf *datadog.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = false
		conf.WithRuntimeMetrics = true
		conf.Tags = datadog.Tags{"tagkey": "tagvalue"}
	})
	assert.NotNil(t, app)
	app.Shutdown()

}

func Test_DataDog_WrapHandler_NotEnabled(t *testing.T) {
	ctx := context.Background()
	app := datadog.NewApplication(ctx, "Glamplify-Unit-Tests", func(conf *datadog.Config) {
		conf.Enabled = false
		conf.ServerlessMode = false
	})

	handler := app.WrapHandler("/", dataDogServeHTTP)
	assert.NotNil(t, handler)
	// assert that handler == dataDogServeHTTP
}

func Test_DataDog_WrapHandler_Enabled(t *testing.T) {
	ctx := context.Background()
	app := datadog.NewApplication(ctx, "Glamplify-Unit-Tests", func(conf *datadog.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = false
	})

	handler := app.WrapHandler("/", dataDogServeHTTP)
	assert.NotNil(t, handler)
	// assert that handler != dataDogServeHTTP
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
	assert.NotNil(t, handler)
}

func dataDogLambadaHandler(ctx context.Context, request aws.ALBTargetGroupRequest) (aws.ALBTargetGroupResponse, error) {
	return aws.ALBTargetGroupResponse{StatusCode: 200}, nil
}

func Test_DataDog_RealWorld_RecordMetric(t *testing.T) {
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
	assert.NotNil(t, span)
}