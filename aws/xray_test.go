package aws

import (
	"context"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New_Tracer(t *testing.T) {
	ctx := context.Background()

	xray := NewTracer(ctx, func(config *TracerConfig) {
		config.EnableLogging = false
		config.Version = "1.0.0"
		config.Environment = "production"
		config.AWSService = "ECS"
	})
	assert.NotNil(t, xray)

	xray = NewTracer(ctx, func(config *TracerConfig) {
		config.EnableLogging = true
		config.Version = "1.0.0"
		config.Environment = "production"
		config.AWSService = "EC2"
	})
	assert.NotNil(t, xray)
}

func Test_Trace_ID(t *testing.T) {
	ctx := context.Background()

	xray := NewTracer(ctx, func(config *TracerConfig) {
		config.EnableLogging = false
		config.Version = "1.0.0"
		config.Environment = "local"
		config.AWSService = "ECS"
	})

	traceID := xray.GetTraceID(ctx)
	assert.Empty(t, traceID)
}

func Test_Segment(t *testing.T) {
	ctx := context.Background()

	xray := NewTracer(ctx, func(config *TracerConfig) {
		config.EnableLogging = false
		config.Version = "1.0.0"
		config.Environment = "local"
		config.AWSService = "ECS"
	})

	h := xray.SegmentHandler("test", mockHandler{})
	assert.NotNil(t, h)

	h = xray.DynamicSegmentHandler("test2", "*", mockHandler{})
	assert.NotNil(t, h)
}

type mockHandler struct {
	mock.Mock
}

func (mh mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}