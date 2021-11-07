package auth

import (
	"context"

	"github.com/cultureamp/glamplify/jwt"
)

type payloadContextKey string

const key = payloadContextKey("payload")

type ValidatedJWTPayload struct {
	Validated bool
	Token     string
	Payload   jwt.Payload
}

func ContextWithValidatedJWTPayload(parent context.Context, payload ValidatedJWTPayload) context.Context {
	ctx := context.WithValue(parent, key, payload)
	return ctx
}

// GetJWTPayload retrieves the authorized user information off the request,
// returning false if the token was not present or failed validation.
func GetJWTPayload(ctx context.Context) (ValidatedJWTPayload, bool) {
	value := ctx.Value(key)

	payload, ok := value.(ValidatedJWTPayload)
	if ok && payload.Validated {
		return payload, true
	}

	return ValidatedJWTPayload{}, false
}

// ContextHasValidatedJWT returns true if the supplied context contains a
// validated JWT payload against the same token as that supplied, indicating
// that the current request's Authorization header has been validated
// successfully.
func ContextHasValidatedJWT(ctx context.Context, token string) bool {
	jwt, ok := GetJWTPayload(ctx)
	if ok &&
		jwt.Validated &&
		jwt.Token == token {

		return true
	}

	return false
}
