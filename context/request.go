package context

import (
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/cultureamp/glamplify/jwt"
	"net/http"
)

const (
	TraceIDHeader = xray.TraceIDHeaderKey // "x-amzn-trace-id"
	RequestIDHeader = "X-Request-ID"
	CorrelationIDHeader = "X-Correlation-ID"
)

func GetRequestScopedFieldsFromRequest(r *http.Request) (RequestScopedFields, bool) {
	return GetRequestScopedFields(r.Context())
}

func AddRequestScopedFieldsRequest(r *http.Request, requestScopeFields RequestScopedFields) *http.Request {
	ctx := AddRequestFields(r.Context(), requestScopeFields)
	return r.WithContext(ctx)
}

// WrapRequest returns the same *http.Request if RequestScopedFields is already present in the context.
// If missing, then it checks http.Request Headers for TraceID, RequestID, and CorrelationID.
// Then this method also tries to decode the JWT payload and adds CustomerAggregateID and UserAggregateID if successful.
func WrapRequest(r *http.Request) *http.Request {

	jwt, err := jwt.NewDecoder() // reads AUTH_PUBLIC_KEY environment var - use PayloadFromRequest() if you want a custom decoder
	if err != nil {
		// TODO - how to log this error? Does it really matter?
	}
	return WrapRequestWithDecoder(r, jwt)
}

// WrapRequestWithDecoder returns the same *http.Request if RequestScopedFields is already present in the context.
// If missing, then it checks http.Request Headers for TraceID, RequestID, and CorrelationID.
// Then this method also tries to decode the JWT payload and adds CustomerAggregateID and UserAggregateID if successful.
func WrapRequestWithDecoder(r *http.Request, jwtDecoder jwt.DecodeJwtToken) *http.Request {
	rsFields, ok := GetRequestScopedFieldsFromRequest(r)
	if ok {
		return r
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
	return r.WithContext(ctx)
}

