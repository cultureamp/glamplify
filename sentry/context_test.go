package sentry_test

import (
	"context"
	"testing"

	"github.com/cultureamp/glamplify/sentry"
	"github.com/stretchr/testify/assert"
)

func TestContext_Fail(t *testing.T) {

	ctx := context.Background()
	txn, err := sentry.FromContext(ctx)

	assert.Nil(t, txn)
	assert.NotNil(t, err)
}
