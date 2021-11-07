package tracing

import (
	"glamplify/server/datadog"
	"net/http"

	"github.com/cultureamp/glamplify/context"
)

// WrapHTTPClient propagates Glamplify, X-ray and Datadog
// tracing values to the request from the context.
func WrapHTTPClient(client *http.Client) *http.Client {
	if client.Transport == nil {
		client.Transport = http.DefaultTransport
	}

	client.Transport = &tracingRoundTripper{
		base: client.Transport,
	}

	return datadog.WrapHTTPClient(client)
}

type tracingRoundTripper struct {
	base http.RoundTripper
}

func (rt *tracingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {

	scopedFields, fieldsPresent := context.GetRequestScopedFieldsFromRequest(req)

	if fieldsPresent {
		addHeader(req, context.TraceIDHeader, scopedFields.TraceID)
		addHeader(req, context.RequestIDHeader, scopedFields.RequestID)
		addHeader(req, context.CorrelationIDHeader, scopedFields.CorrelationID)
	}

	return rt.base.RoundTrip(req)
}

func addHeader(req *http.Request, headerName string, headerValue string) {
	if headerValue != "" {
		req.Header.Add(headerName, headerValue)
	}
}
