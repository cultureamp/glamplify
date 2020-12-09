package context

import (
	"context"
)

// RequestScopedFields represents scoping information
type RequestScopedFields struct {
	TraceID             string `json:"trace_id"`		// AWS XRAY trace id. Format of this is controlled by AWS. Do not rely on it, some services may not use XRAY.
	RequestID           string `json:"request_id"`		// Client generated RANDOM string. Most of the time this will be empty. Clients can set this to help us diagnose issues.
	CorrelationID       string `json:"correlation_id"`	// Set ALWAYS by the web-gateway as a UUID v4.
	UserAggregateID     string `json:"user"`			// If JWT and correct key present, then this will be set to the Effective User UUID
	CustomerAggregateID string `json:"customer"`		// If JWT and correct key present, then this will be set to the Customer UUID (aka Account)
}

// NewRequestScopeFields creates a new RequestScopedFields
func NewRequestScopeFields(traceID string, requestID string, correlationID string, customer string, user string) RequestScopedFields {
	return RequestScopedFields{
		TraceID:             traceID,
		RequestID:           requestID,
		CorrelationID:       correlationID,
		CustomerAggregateID: customer,
		UserAggregateID:     user,
	}
}

// AddToCtx add a RequestScopedFields to the context
func (rsFields RequestScopedFields) AddToCtx(ctx context.Context) context.Context {
	return AddRequestFields(ctx, rsFields)
}
