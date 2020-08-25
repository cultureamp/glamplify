package newrelic_test

import (
	"context"
	"errors"
	"github.com/cultureamp/glamplify/log"
	"github.com/cultureamp/glamplify/newrelic"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/assert"
)

func TestTxn_AddAttribute_Server_Success(t *testing.T) {
	app, err := newrelic.NewApplication("Glamplify-Unit-Tests", func(conf *newrelic.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = false
	})

	assert.Assert(t, err == nil, err)
	assert.Assert(t, app != nil, "application was nil")

	_, handler := app.WrapHTTPHandler("/", addAttribute)
	h := http.HandlerFunc(handler)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(rr, req)

	app.Shutdown()
}

func addAttribute(w http.ResponseWriter, r *http.Request) {
	txn, err := newrelic.TxnFromRequest(w, r)
	if err == nil {
		txn.AddAttributes(log.Fields{
			"aString": "hello world",
			"aInt":    123,
		})
	}
}

func TestTxn_NoticeError_Server_Success(t *testing.T) {
	app, err := newrelic.NewApplication("Glamplify-Unit-Tests", func(conf *newrelic.Config) {
		conf.Enabled = true
		conf.Logging = true
		conf.ServerlessMode = false
	})

	assert.Assert(t, err == nil, err)
	assert.Assert(t, app != nil, "application was nil")

	_, handler := app.WrapHTTPHandler("/", reportError)
	h := http.HandlerFunc(handler)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)

	// Add *testing.T to request context
	ctx := req.Context()
	ctx = context.WithValue(ctx, "t", t)
	req = req.WithContext(ctx)
	h.ServeHTTP(rr, req)

	app.Shutdown()
}

func reportError(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	t, _ := ctx.Value("t").(*testing.T)

	txn, err := newrelic.TxnFromRequest(w, r)
	if err == nil {
		err = txn.ReportError(errors.New("standard error message"))
		assert.Assert(t, err == nil, err )
		txn.ReportErrorDetails("detailed error", "txn_test", log.Fields{
			"aString": "hello world",
		})
		assert.Assert(t, err == nil, err )
	}
}

