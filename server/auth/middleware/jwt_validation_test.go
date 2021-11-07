package middleware

import (
	"context"
	"errors"
	"glamplify/server/auth"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cultureamp/glamplify/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMiddlewareValidationOccurs(t *testing.T) {
	token := "jwt-token-header-value"

	payloadResult := jwt.Payload{
		RealUser:      "real",
		EffectiveUser: "effective-user",
		Customer:      "customer",
	}
	decoder := &testDecoder{}
	decoder.On("Decode", token).Return(payloadResult, nil)

	var actualContext context.Context
	nextHandler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		actualContext = r.Context()
	})

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer "+token)

	sut := NewJWTValidationMiddleware(decoder)(nextHandler)

	sut.ServeHTTP(httptest.NewRecorder(), r)

	expectedPayload := auth.ValidatedJWTPayload{
		Validated: true,
		Token:     token,
		Payload:   payloadResult,
	}

	actualPayload, ok := auth.GetJWTPayload(actualContext)
	require.True(t, ok)

	assert.Equal(t, expectedPayload, actualPayload)
	decoder.AssertExpectations(t)
}

func TestBearerValidateTokenSucceeds(t *testing.T) {
	token := "jwt-token-supplied"

	payloadResult := jwt.Payload{
		RealUser:      "real",
		EffectiveUser: "effective-user",
		Customer:      "customer",
	}
	decoder := &testDecoder{}
	decoder.On("Decode", token).Return(payloadResult, nil)

	sut := jwtValidationMiddleware{
		decoder: decoder,
	}

	payload := sut.validateToken(context.Background(), token)

	expectedPayload := auth.ValidatedJWTPayload{
		Validated: true,
		Token:     token,
		Payload:   payloadResult,
	}

	assert.Equal(t, expectedPayload, payload)
	decoder.AssertExpectations(t)
}

func TestBearerValidateTokenFails(t *testing.T) {
	token := "jwt-token-supplied"

	payloadResult := jwt.Payload{}
	decodeError := errors.New("decode failed")

	decoder := &testDecoder{}
	decoder.On("Decode", token).Return(payloadResult, decodeError)

	sut := jwtValidationMiddleware{
		decoder: decoder,
	}

	payload := sut.validateToken(context.Background(), token)

	expectedPayload := auth.ValidatedJWTPayload{
		Validated: false,
		Token:     token,
		Payload:   payloadResult,
	}

	assert.Equal(t, expectedPayload, payload)
	decoder.AssertExpectations(t)
}

func TestBearerTokenPresent(t *testing.T) {
	headerValue := "Bearer foo"
	token := getBearerToken(headerValue)

	expected := "foo"
	assert.Equal(t, expected, token)
}

func TestBearerTokenInvalidOrNotPresent(t *testing.T) {
	cases := []string{"", "BearerFoo", "ABC", "Bearer", "Bearer ", "bearer foo"}

	for _, headerValue := range cases {
		t.Run(headerValue, func(t *testing.T) {
			token := getBearerToken(headerValue)
			expected := ""
			assert.Equal(t, expected, token)
		})
	}
}

type testDecoder struct {
	mock.Mock
}

func (d *testDecoder) Decode(tokenString string) (jwt.Payload, error) {
	args := d.Called(tokenString)
	return args.Get(0).(jwt.Payload), args.Error(1)
}
