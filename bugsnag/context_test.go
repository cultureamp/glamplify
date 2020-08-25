package bugsnag_test

import (
	"context"
	"github.com/cultureamp/glamplify/bugsnag"
	"gotest.tools/assert"
	"testing"
)

func TestContext_Fail(t *testing.T) {

	ctx := context.TODO()
	txn, err := bugsnag.BugsnagFromContext(ctx)

	assert.Assert(t, txn == nil, txn)
	assert.Assert(t, err != nil, err)
}