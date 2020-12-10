package sentry_test

import (
	"context"
	"testing"

	"github.com/cultureamp/glamplify/sentry"
	"gotest.tools/assert"
)

func TestContext_Fail(t *testing.T) {

	ctx := context.Background()
	txn, err := sentry.FromContext(ctx)

	assert.Assert(t, txn == nil, txn)
	assert.Assert(t, err != nil, err)
}
