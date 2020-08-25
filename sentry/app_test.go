package sentry_test

import (
	"context"
	"errors"
	"github.com/cultureamp/glamplify/sentry"
	"gotest.tools/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSentry_Error_Success(t *testing.T) {

	ctx := context.Background()
	sentry, err := sentry.NewApplication(ctx, "GlamplifyUnitTests", func (conf *sentry.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.AppVersion = "1.0.0"
	})
	assert.Assert(t, err == nil, err)

	id := sentry.Error(errors.New("NPE"))
	assert.Assert(t, id != nil, id)

	sentry.Shutdown()
}

func TestSentry_Context_Success(t *testing.T) {

	ctx := context.Background()
	sentry, err := sentry.NewApplication(ctx, "GlamplifyUnitTests", func (conf *sentry.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.AppVersion = "1.0.0"
	})
	assert.Assert(t, err == nil, err)

	_, handler := sentry.WrapHTTPHandler("/", rootRequest)
	h := http.HandlerFunc(handler)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)

	// Add *testing.T to request context
	ctx = req.Context()
	ctx = context.WithValue(ctx, "t", t)
	req = req.WithContext(ctx)

	h.ServeHTTP(rr, req)

	sentry.Shutdown()
}

func rootRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	t, _ := ctx.Value("t").(*testing.T)

	sentry, err := sentry.SentryFromContext(ctx)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, sentry != nil, sentry)
}


