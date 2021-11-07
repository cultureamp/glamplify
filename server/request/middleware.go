package request

import (
	"context"
	"glamplify/server/auth"
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
	gcontext "github.com/cultureamp/glamplify/context"
	"github.com/google/uuid"
	"goa.design/goa/v3/middleware"
)

// NewRequestTracingContextMiddleware returns request middleware that ensures
// that the request context contains the expected values for cross-service
// tracing, as well as details of the currently authorized user.
//
// Tracing details are assigned based on request header values, or generated
// when the header values are not present.
func NewRequestTracingContextMiddleware() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := contextAddRequestScopedFields(r)
			ctx = contextAddGoaRequestID(ctx)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func contextAddRequestScopedFields(r *http.Request) context.Context {
	ctx := r.Context()

	// cross-service tracing context from incoming request headers

	traceID := r.Header.Get(gcontext.TraceIDHeader)
	if xray.RequestWasTraced(ctx) {
		traceID = xray.TraceID(ctx)
	}

	requestID := r.Header.Get(gcontext.RequestIDHeader)
	if requestID == "" {
		requestID = uuid.New().String()
	}

	correlationID := r.Header.Get(gcontext.CorrelationIDHeader)
	if correlationID == "" {
		correlationID = uuid.New().String()
	}

	// current request authenticated user context

	customerAggregateID := ""
	userAggregateID := ""
	if jwt, ok := auth.GetJWTPayload(ctx); ok {
		customerAggregateID = jwt.Payload.Customer
		userAggregateID = jwt.Payload.EffectiveUser
	}

	rsFields := gcontext.NewRequestScopeFields(traceID, requestID, correlationID, customerAggregateID, userAggregateID)

	ctx = rsFields.AddToCtx(ctx)

	return ctx
}

// contextAddGoaRequestID uses the Glamplify request ID for the Goa request ID,
// adding compatibility with Goa's error types
func contextAddGoaRequestID(ctx context.Context) context.Context {
	if rs, ok := gcontext.GetRequestScopedFields(ctx); ok {
		ctx = context.WithValue(ctx, middleware.RequestIDKey, rs.RequestID)
	}
	return ctx
}
