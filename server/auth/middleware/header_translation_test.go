package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthHeaderInjection(t *testing.T) {
	// create a handler to use as "next" which will verify the request
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// here assert that the authorization header has been set
		assert.Equal(t, r.Header.Get("Authorization"), "token")
	})
	// create the handler to test, using our custom "next" handler
	handlerToTest := AuthHeaderTranslator()(nextHandler)

	// create a mock request to use
	req := httptest.NewRequest("GET", "http://testing", nil)

	// set header on request
	req.Header.Add(BFFCustomAuthHeader, "token")

	// call the handler using a mock response recorder (we'll not use that anyway)
	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
}
