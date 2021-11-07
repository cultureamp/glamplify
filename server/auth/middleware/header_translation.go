package middleware

import "net/http"

const BFFCustomAuthHeader = "X-CA-SGW-Authorization"

// AuthHeaderTranslator provides middleware that translates the tunneled
// "X-CA-SGW-Authorization" header to the "Authorization" header. This behaviour
// is required when an IAM authorizer is in place at the API Gateway level,
// which makes use of the Authorization header.
func AuthHeaderTranslator() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get(BFFCustomAuthHeader)
			if authHeader != "" {
				r.Header.Set("Authorization", authHeader)
			}
			h.ServeHTTP(w, r)
		})
	}
}
