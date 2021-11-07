package auth

import (
	"context"
	"testing"

	"github.com/cultureamp/glamplify/jwt"
	"github.com/stretchr/testify/assert"
)

func TestContextWithJWTPayload(t *testing.T) {
	ctx := context.Background()
	payload := ValidatedJWTPayload{
		Validated: true,
		Token:     "token value",
		Payload: jwt.Payload{
			RealUser:      "real user",
			EffectiveUser: "eff user",
			Customer:      "customer",
		},
	}

	ctx = ContextWithValidatedJWTPayload(ctx, payload)

	assert.Equal(t, ctx.Value(key), payload)
}

func TestGetJWTPayloadMissing(t *testing.T) {
	ctx := context.Background()

	_, ok := GetJWTPayload(ctx)
	assert.False(t, ok)
}

func TestGetJWTPayloadValid(t *testing.T) {
	ctx := context.Background()
	payload := ValidatedJWTPayload{
		Validated: true,
		Token:     "token value",
		Payload: jwt.Payload{
			RealUser:      "real user",
			EffectiveUser: "eff user",
			Customer:      "customer",
		},
	}

	ctx = ContextWithValidatedJWTPayload(ctx, payload)

	v, ok := GetJWTPayload(ctx)

	assert.True(t, ok)
	assert.Equal(t, payload, v)
}

func TestGetJWTPayloadInvalid(t *testing.T) {
	ctx := context.Background()
	payload := ValidatedJWTPayload{
		Validated: false,
		Token:     "token value",
		Payload: jwt.Payload{
			RealUser:      "real user",
			EffectiveUser: "eff user",
			Customer:      "customer",
		},
	}

	ctx = ContextWithValidatedJWTPayload(ctx, payload)

	v, ok := GetJWTPayload(ctx)

	assert.False(t, ok)
	assert.Equal(t, "", v.Token)
	assert.Equal(t, jwt.Payload{}, v.Payload)
}

func TestContextHasValidJWTSucceeds(t *testing.T) {
	payload := ValidatedJWTPayload{
		Validated: true,
		Token:     "token value",
		Payload: jwt.Payload{
			RealUser:      "real user",
			EffectiveUser: "eff user",
			Customer:      "customer",
		},
	}
	ctx := ContextWithValidatedJWTPayload(context.Background(), payload)

	ok := ContextHasValidatedJWT(ctx, payload.Token)

	assert.True(t, ok)
}

func TestContextHasValidJWTFailsWhenNotPresent(t *testing.T) {
	ok := ContextHasValidatedJWT(context.Background(), "")
	assert.False(t, ok)
}

func TestContextHasValidJWTFailsWhenNotValid(t *testing.T) {
	payload := ValidatedJWTPayload{
		Validated: false,
		Token:     "token value",
		Payload:   jwt.Payload{},
	}
	ctx := ContextWithValidatedJWTPayload(context.Background(), payload)

	ok := ContextHasValidatedJWT(ctx, payload.Token)

	assert.False(t, ok)
}

func TestContextHasValidJWTFailsWithTokenMismatch(t *testing.T) {
	payload := ValidatedJWTPayload{
		Validated: true,
		Token:     "token value",
		Payload:   jwt.Payload{},
	}
	ctx := ContextWithValidatedJWTPayload(context.Background(), payload)

	ok := ContextHasValidatedJWT(ctx, payload.Token+"yikes")

	assert.False(t, ok)
}
