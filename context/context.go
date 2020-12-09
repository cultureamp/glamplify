package context

import (
	"context"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/google/uuid"
)

// AddRequestFields adds a RequestScopedFields to the context
func AddRequestFields(ctx context.Context, rsFields RequestScopedFields) context.Context {
	return context.WithValue(ctx, RequestFieldsCtx, rsFields)
}

// GetRequestScopedFields gets the RequestScopedFields from the context
func GetRequestScopedFields(ctx context.Context) (RequestScopedFields, bool) {
	rsFields, ok := ctx.Value(RequestFieldsCtx).(RequestScopedFields)
	return rsFields, ok
}

// WrapCtx initializes a context with default RequestScopedFields
func WrapCtx(ctx context.Context) context.Context {

	rsFields, ok := GetRequestScopedFields(ctx)
	if ok {
		return ctx
	}

	traceID := ""
	if xray.RequestWasTraced(ctx) {
		traceID = xray.TraceID(ctx)
	}

	requestID := uuid.New().String()
	correlationID := uuid.New().String()

	rsFields = RequestScopedFields{
		TraceID: traceID,
		RequestID: requestID,
		CorrelationID: correlationID,
	}

	return AddRequestFields(ctx, rsFields)
}
