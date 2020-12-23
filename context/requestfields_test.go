package context_test

import (
	"context"
	"testing"

	context2 "github.com/cultureamp/glamplify/context"
	"github.com/stretchr/testify/assert"
)

func Test_TransactionFields_New(t *testing.T) {
	transactionFields := context2.NewRequestScopeFields("1-2-3", "7-8-9", "1-5-9", "hooli", "UserAggregateID-123")
	assert.Equal(t, "1-2-3",  transactionFields.TraceID)
	assert.Equal(t, "7-8-9",  transactionFields.RequestID)
	assert.Equal(t, "1-5-9",  transactionFields.CorrelationID)
	assert.Equal(t, "hooli",  transactionFields.CustomerAggregateID)
	assert.Equal(t, "UserAggregateID-123",  transactionFields.UserAggregateID)
}

func Test_TransactionFields_NewFromCtx(t *testing.T) {
	transactionFields := context2.NewRequestScopeFields("1-2-3", "7-8-9","1-5-9", "hooli", "UserAggregateID-123")

	ctx := context.Background()
	ctx = transactionFields.AddToCtx(ctx)

	rsFields, ok := context2.GetRequestScopedFields(ctx)
	assert.True(t, ok)
	assert.Equal(t, "1-2-3",  transactionFields.TraceID)
	assert.Equal(t, "7-8-9",  transactionFields.RequestID)
	assert.Equal(t, "1-5-9",  rsFields.CorrelationID)
	assert.Equal(t, "hooli",  transactionFields.CustomerAggregateID)
	assert.Equal(t, "UserAggregateID-123",  transactionFields.UserAggregateID)
}

