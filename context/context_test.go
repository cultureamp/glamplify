package context_test

import (
	"context"
	"testing"

	gcontext "github.com/cultureamp/glamplify/context"
	"github.com/stretchr/testify/assert"
)

func Test_Context_AddGet(t *testing.T) {

	ctx := context.Background()
	ctx = gcontext.AddRequestFields(ctx, gcontext.RequestScopedFields{
		TraceID:             "trace1",
		RequestID:           "request1",
		CorrelationID:       "correlation1",
		CustomerAggregateID: "cust1",
		UserAggregateID:     "user1",
	})

	rsFields, ok := gcontext.GetRequestScopedFields(ctx)
	assert.True(t, ok)
	assert.Equal(t, "trace1", rsFields.TraceID)
	assert.Equal(t, "request1", rsFields.RequestID)
	assert.Equal(t, "correlation1", rsFields.CorrelationID)
	assert.Equal(t, "cust1", rsFields.CustomerAggregateID)
	assert.Equal(t, "user1", rsFields.UserAggregateID)
}

func Test_Context_TraceID_AddGet_Empty(t *testing.T) {

	ctx := context.Background()
	ctx = gcontext.AddRequestFields(ctx, gcontext.RequestScopedFields{
	})

	rsFields, ok := gcontext.GetRequestScopedFields(ctx)
	assert.True(t, ok)
	assert.Empty(t, rsFields.TraceID)
	assert.Empty(t, rsFields.RequestID)
	assert.Empty(t, rsFields.CorrelationID)
	assert.Empty(t, rsFields.CustomerAggregateID)
	assert.Empty(t, rsFields.UserAggregateID)
}

func Test_Context_Wrap(t *testing.T) {

	ctx := context.Background()
	ctx = gcontext.WrapCtx(ctx)

	rsFields, ok := gcontext.GetRequestScopedFields(ctx)
	assert.True(t, ok)
	assert.Empty(t, rsFields.TraceID)
	assert.NotEmpty(t, rsFields.RequestID)
	assert.NotEmpty(t, rsFields.CorrelationID)
	assert.Empty(t, rsFields.CustomerAggregateID)
	assert.Empty(t, rsFields.UserAggregateID)
}