package sentry_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/cultureamp/glamplify/sentry"
	sentrygo "github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
)

func Test_Sentry_New(t *testing.T) {
	ctx := context.Background()

	app, err := sentry.NewApplication(ctx, "GlamplifyUnitTests", func (conf *sentry.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.AppVersion = "1.0.0"
		conf.DSN = "https://177fbd4b35304a80aeaef835f938de69@o19604.ingest.app.io/5447011"
	})
	assert.Nil(t, err)
	assert.NotNil(t, app)


	app, err = sentry.NewApplication(ctx, "GlamplifyUnitTests", func (conf *sentry.Config) {
		conf.Enabled = false
		conf.Logging = false
		conf.AppVersion = "1.0.0"
		conf.DSN = "https://177fbd4b35304a80aeaef835f938de69@o19604.ingest.app.io/5447011"
	})
	assert.Nil(t, err)
	assert.NotNil(t, app)
}

func Test_Sentry_Error_Success(t *testing.T) {

	ctx := context.Background()
	sentry, err := sentry.NewApplication(ctx, "GlamplifyUnitTests", func (conf *sentry.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.AppVersion = "1.0.0"
		conf.DSN = "https://177fbd4b35304a80aeaef835f938de69@o19604.ingest.sentry.io/5447011"
	})
	assert.Nil(t, err)

	id := sentry.Error(errors.New("glamplify test NPE"))
	assert.NotNil(t, id)

	sentry.Shutdown()
}

func Test_Sentry_Context_Success(t *testing.T) {

	ctx := context.Background()
	sentry, err := sentry.NewApplication(ctx, "GlamplifyUnitTests", func (conf *sentry.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.AppVersion = "1.0.0"
		conf.DSN = "https://177fbd4b35304a80aeaef835f938de69@o19604.ingest.sentry.io/5447011"
	})
	assert.Nil(t, err)

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

func Test_Sentry_Sync_Transport(t *testing.T) {

	ctx := context.Background()
	sentry, err := sentry.NewApplication(ctx, "GlamplifyUnitTests", func (conf *sentry.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.AppVersion = "1.0.0"
		conf.DSN = "https://177fbd4b35304a80aeaef835f938de69@o19604.ingest.sentry.io/5447011"
		conf.Transport = &sentrygo.HTTPSyncTransport{Timeout: 1500 * time.Millisecond}
	})
	assert.Nil(t, err)

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

func Test_Sentry_WrapHandler(t *testing.T) {
	ctx := context.Background()
	sentry, _ := sentry.NewApplication(ctx, "GlamplifyUnitTests", func (conf *sentry.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.AppVersion = "1.0.0"
		conf.DSN = "https://177fbd4b35304a80aeaef835f938de69@o19604.ingest.sentry.io/5447011"
	})

	pattern, handler := sentry.WrapHTTPHandler("/", rootRequest)
	assert.NotNil(t, handler)
	assert.Equal(t, "/", pattern)
}

func rootRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	t, _ := ctx.Value("t").(*testing.T)

	sentry, err := sentry.FromContext(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, sentry)

	id := sentry.Message("glamplify http handler test message")
	assert.NotEmpty(t, id)
}
