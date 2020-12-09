package context

import (
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/cultureamp/glamplify/jwt"
	"net/http"
)

const (
	// TraceIDHeader = xray.TraceIDHeaderKey eg. "x-amzn-trace-id"
	TraceIDHeader = xray.TraceIDHeaderKey
	// RequestIDHeader = "X-Request-ID"
	RequestIDHeader = "X-Request-ID"
	// CorrelationIDHeader = "X-Correlation-ID"
	CorrelationIDHeader = "X-Correlation-ID"
	// ErrorUUID = "00000000-0000-0000-0000-000000000000"
	ErrorUUID = "00000000-0000-0000-0000-000000000000"
)

// GetRequestScopedFieldsFromRequest gets the RequestScopedFields from the request context
func GetRequestScopedFieldsFromRequest(r *http.Request) (RequestScopedFields, bool) {
	return GetRequestScopedFields(r.Context())
}

// AddRequestScopedFieldsRequest adds a RequestScopedFields to the request context
func AddRequestScopedFieldsRequest(r *http.Request, requestScopeFields RequestScopedFields) *http.Request {
	ctx := AddRequestFields(r.Context(), requestScopeFields)
	return r.WithContext(ctx)
}

// WrapRequest returns the same *http.Request if RequestScopedFields is already present in the context.
// If missing, then it checks http.Request Headers for TraceID, RequestID, and CorrelationID.
// Then this method also tries to decode the JWT payload and adds CustomerAggregateID and UserAggregateID if successful.
func WrapRequest(r *http.Request) (*http.Request, error) {

	// reads AUTH_PUBLIC_KEY environment var - use PayloadFromRequest() if you want a custom decoder
	// then use WrapRequestWithDecoder
	jwt, err := jwt.NewDecoder()
	if err != nil {
		return r, err
	}
	return WrapRequestWithDecoder(r, jwt)
}

// WrapRequestWithDecoder returns the same *http.Request if RequestScopedFields is already present in the context.
// If missing, then it checks http.Request Headers for TraceID, RequestID, and CorrelationID.
// Then this method also tries to decode the JWT payload and adds CustomerAggregateID and UserAggregateID if successful.
func WrapRequestWithDecoder(r *http.Request, jwtDecoder jwt.DecodeJwtToken) (*http.Request, error) {
	rsFields, ok := GetRequestScopedFieldsFromRequest(r)
	if ok {
		return r, nil
	}

	// need to create new RequestScopedFields
	ctx := r.Context()
	traceID := r.Header.Get(TraceIDHeader)
	requestID := r.Header.Get(RequestIDHeader)
	correlationID := r.Header.Get(CorrelationIDHeader)

	payload, err := jwt.PayloadFromRequest(r, jwtDecoder)

	if err == nil {
		rsFields = NewRequestScopeFields(traceID, requestID, correlationID, payload.Customer, payload.EffectiveUser)
	} else {
		rsFields = NewRequestScopeFields(traceID, requestID, correlationID, "", "")
	}

	ctx = rsFields.AddToCtx(ctx)
	return r.WithContext(ctx), err
}

