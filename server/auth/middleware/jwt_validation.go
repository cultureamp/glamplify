package middleware

import (
	"context"
	"glamplify/server/auth"
	"net/http"
	"strings"

	"github.com/cultureamp/glamplify/jwt"
	"github.com/cultureamp/glamplify/log"
)

// Decoder describes the contract required to decode a JWT token
type Decoder interface {
	// Decode takes a JWT token and verifies it, returning the decoded payload
	// if validation is successful or an error if the token fails.
	Decode(tokenString string) (jwt.Payload, error)
}

type jwtValidationMiddleware struct {
	next    http.Handler
	decoder Decoder
}

// NewJWTValidationMiddleware supplies middleware that will decode a JWT present
// in the Authorization header, placing the result of this validation on the
// context.
//
// It does *NOT* otherwise modify the request: taking action as a result of
// failed validation is delegated to the handler for the current route (which
// may not require credentials). When using Goa, the generated "Auther" can
// interrogate the context for the validation result and payload.
//
// Details of the decoded JWT is only placed in the context if validation
// succeeds.
func NewJWTValidationMiddleware(decoder Decoder) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		v := jwtValidationMiddleware{
			next:    next,
			decoder: decoder,
		}

		return v
	}
}

func (m jwtValidationMiddleware) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	token := getBearerToken(req.Header.Get("Authorization"))

	v := m.validateToken(ctx, token)
	ctx = auth.ContextWithValidatedJWTPayload(ctx, v)

	m.next.ServeHTTP(resp, req.WithContext(ctx))
}

// validateToken uses the decoder to validate the supplied token, returning the
// resulting payload. No decoded JWT details are returned if it fails
// validation.
func (m jwtValidationMiddleware) validateToken(ctx context.Context, token string) auth.ValidatedJWTPayload {
	v := auth.ValidatedJWTPayload{
		Token: token,
	}

	payload, err := m.decoder.Decode(token)
	if err == nil {
		v.Validated = true
		v.Payload = payload
	} else {
		logger := log.NewFromCtx(ctx)
		logger.Error("jwt_validation_failed", err)
	}

	return v
}

// getBearerToken strips the required "Bearer " prefix from an Authorization header
// value, returning an empty string if the prefix is not present.
func getBearerToken(header string) string {
	if !strings.HasPrefix(header, "Bearer ") {
		return ""
	}

	return header[7:]
}
