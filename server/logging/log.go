package logging

import (
	"context"
	"glamplify/server/auth"
	"net/http"

	"github.com/cultureamp/glamplify/log"
	goa "goa.design/goa/v3/pkg"
)

// NewEndpointMiddleware supplies middleware logs an error resulting from endpoint execution
func NewEndpointMiddleware() func(goa.Endpoint) goa.Endpoint {
	return func(next goa.Endpoint) goa.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			res, err := next(ctx, request)

			if err != nil {
				logger := log.NewFromCtx(ctx)
				logger.Error("request_error", err)
			}

			return res, err
		}
	}
}

// NewRequestMiddleware returns request middleware that extracts some headers from the request and adds them
// to the context for the logger to utilize. It then logs the request.
func NewRequestMiddleware() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)

			ctx := r.Context()

			// datadog adds this value to the context
			coldStart := ctx.Value("cold_start")

			fields := log.Fields{
				"method":     r.Method,
				"path":       r.URL.Path,
				"query":      r.URL.Query().Encode(),
				"cold_start": coldStart,
			}

			if jwt, ok := auth.GetJWTPayload(ctx); ok {
				fields = fields.Merge(log.Fields{
					"customer": jwt.Payload.Customer,
					"user":     jwt.Payload.EffectiveUser,
				})
			}

			logger := log.NewFromCtx(ctx)
			logger.Info("request", fields)
		})
	}
}
