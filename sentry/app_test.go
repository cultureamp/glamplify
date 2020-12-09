package sentry_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cultureamp/glamplify/sentry"
	sentrygo "github.com/getsentry/sentry-go"
	"gotest.tools/assert"
)

func TestSentry_Error_Success(t *testing.T) {

	ctx := context.Background()
	sentry, err := sentry.NewApplication(ctx, "GlamplifyUnitTests", func (conf *sentry.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.AppVersion = "1.0.0"
		conf.DSN = "https://177fbd4b35304a80aeaef835f938de69@o19604.ingest.sentry.io/5447011"
	})
	assert.Assert(t, err == nil, err)

	id := sentry.Error(errors.New("glamplify test NPE"))
	assert.Assert(t, id != nil, id)

	sentry.Shutdown()
}

func TestSentry_Context_Success(t *testing.T) {

	ctx := context.Background()
	sentry, err := sentry.NewApplication(ctx, "GlamplifyUnitTests", func (conf *sentry.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.AppVersion = "1.0.0"
		conf.DSN = "https://177fbd4b35304a80aeaef835f938de69@o19604.ingest.sentry.io/5447011"
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

func TestSentry_Sync_Transport(t *testing.T) {

	ctx := context.Background()
	sentry, err := sentry.NewApplication(ctx, "GlamplifyUnitTests", func (conf *sentry.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.AppVersion = "1.0.0"
		conf.DSN = "https://177fbd4b35304a80aeaef835f938de69@o19604.ingest.sentry.io/5447011"
		conf.Transport = &sentrygo.HTTPSyncTransport{Timeout: 1500 * time.Millisecond}
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

	sentry.Flush(50 * time.Millisecond)
}

func rootRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	t, _ := ctx.Value("t").(*testing.T)

	sentry, err := sentry.SentryFromContext(ctx)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, sentry != nil, sentry)

	id := sentry.Message("glamplify http handler test message")
	assert.Assert(t, *id != "", id)
}


