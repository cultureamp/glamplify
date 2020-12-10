package aws

import (
	"context"
	"testing"

	"gotest.tools/assert"
)

func Test_New_Tracer(t *testing.T) {
	ctx := context.Background()

	xray := NewTracer(ctx, func(config *TracerConfig) {
		config.EnableLogging = false
		config.Version = "1.0.0"
		config.Environment = "local"
		config.AWSService = "ECS"
	})
	assert.Assert(t, xray != nil, xray)
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
	assert.Assert(t, traceID == "", traceID)
}