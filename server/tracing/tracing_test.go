package tracing

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	gcontext "github.com/cultureamp/glamplify/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTracingHeadersPresent(t *testing.T) {
	ctx := context.Background()
	rsFields := gcontext.NewRequestScopeFields("trace-id", "request-id", "correlation-id", "", "")
	ctx = rsFields.AddToCtx(ctx)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://www.example.com/happy_path", nil)
	w := httptest.NewRecorder()

	traceHeaderValue := ""
	requestIDHeaderValue := ""
	corrIDHeaderValue := ""

	trippy := &testRoundTripper{
		impl: func(req *http.Request) (*http.Response, error) {
			traceHeaderValue = req.Header.Get(gcontext.TraceIDHeader)
			requestIDHeaderValue = req.Header.Get(gcontext.RequestIDHeader)
			corrIDHeaderValue = req.Header.Get(gcontext.CorrelationIDHeader)

			w.WriteHeader(http.StatusTeapot)

			return w.Result(), nil
		},
	}

	sut := &tracingRoundTripper{
		base: trippy,
	}

	response, err := sut.RoundTrip(req)
	response.Body.Close()

	require.Nil(t, err)
	assert.Equal(t, rsFields.TraceID, traceHeaderValue)
	assert.Equal(t, rsFields.RequestID, requestIDHeaderValue)
	assert.Equal(t, rsFields.CorrelationID, corrIDHeaderValue)

	assert.Equal(t, http.StatusTeapot, response.StatusCode)
}

func TestTracingHeadersMissingWhenContextEmpty(t *testing.T) {
	ctx := context.Background()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://www.example.com/happy_path", nil)
	w := httptest.NewRecorder()

	traceHeaderValue := ""
	requestIDHeaderValue := ""
	corrIDHeaderValue := ""

	trippy := &testRoundTripper{
		impl: func(req *http.Request) (*http.Response, error) {
			traceHeaderValue = req.Header.Get(gcontext.TraceIDHeader)
			requestIDHeaderValue = req.Header.Get(gcontext.RequestIDHeader)
			corrIDHeaderValue = req.Header.Get(gcontext.CorrelationIDHeader)

			w.WriteHeader(http.StatusTeapot)

			return w.Result(), nil
		},
	}
	sut := &tracingRoundTripper{
		base: trippy,
	}

	response, err := sut.RoundTrip(req)
	response.Body.Close()

	require.Nil(t, err)
	assert.Equal(t, "", traceHeaderValue)
	assert.Equal(t, "", requestIDHeaderValue)
	assert.Equal(t, "", corrIDHeaderValue)

	assert.Equal(t, http.StatusTeapot, response.StatusCode)
}

type testRoundTripper struct {
	impl func(*http.Request) (*http.Response, error)
}

func (rt *testRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	return rt.impl(req)
}
