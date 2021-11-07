package request

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	gcontext "github.com/cultureamp/glamplify/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestScopedFieldsAddedToContext(t *testing.T) {
	var ctx context.Context

	// create a handler to use as "next" which will verify the request
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	})

	// create the handler to test, using our custom "next" handler
	handlerToTest := NewRequestTracingContextMiddleware()(nextHandler)

	// create a mock request to use
	req := httptest.NewRequest("GET", "http://testing", nil)

	// set header on request
	req.Header.Add("x-amzn-trace-id", "traceID-123")
	req.Header.Add("X-Request-ID", "requestID-123")
	req.Header.Add("X-Correlation-ID", "correlationID-123")

	// call the handler using a mock response recorder (we'll not use that anyway)
	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)

	rs, ok := gcontext.GetRequestScopedFields(ctx)

	require.True(t, ok)
	assert.Equal(t, "traceID-123", rs.TraceID)
	assert.Equal(t, "requestID-123", rs.RequestID)
	assert.Equal(t, "correlationID-123", rs.CorrelationID)
}

func TestWhenRequestScopedFieldsAreMissing(t *testing.T) {
	var ctx context.Context

	// create a handler to use as "next" which will verify the request
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
	})

	// create the handler to test, using our custom "next" handler
	handlerToTest := NewRequestTracingContextMiddleware()(nextHandler)

	// create a mock request to use
	req := httptest.NewRequest("GET", "http://testing", nil)

	// call the handler using a mock response recorder (we'll not use that anyway)
	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)

	rs, ok := gcontext.GetRequestScopedFields(ctx)

	require.True(t, ok)
	assert.Equal(t, "", rs.TraceID)
	assert.NotEmpty(t, rs.RequestID)
	assert.NotEmpty(t, rs.CorrelationID)
}
