package context_test

import (
	"context"
	gcontext "github.com/cultureamp/glamplify/context"
	"gotest.tools/assert"
	"testing"
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
	assert.Assert(t, ok, ok)
	assert.Assert(t, rsFields.TraceID == "trace1", rsFields)
	assert.Assert(t, rsFields.RequestID == "request1", rsFields)
	assert.Assert(t, rsFields.CorrelationID == "correlation1", rsFields)
	assert.Assert(t, rsFields.CustomerAggregateID == "cust1", rsFields)
	assert.Assert(t, rsFields.UserAggregateID == "user1", rsFields)
}

func Test_Context_TraceID_AddGet_Empty(t *testing.T) {

	ctx := context.Background()
	ctx = gcontext.AddRequestFields(ctx, gcontext.RequestScopedFields{
	})

	rsFields, ok := gcontext.GetRequestScopedFields(ctx)
	assert.Assert(t, ok, ok)
	assert.Assert(t, rsFields.TraceID == "", rsFields)
	assert.Assert(t, rsFields.RequestID == "", rsFields)
	assert.Assert(t, rsFields.CorrelationID == "", rsFields)
	assert.Assert(t, rsFields.CustomerAggregateID == "", rsFields)
	assert.Assert(t, rsFields.UserAggregateID == "", rsFields)
}

func Test_Context_Wrap(t *testing.T) {

	ctx := context.Background()
	ctx = gcontext.WrapCtx(ctx)

	rsFields, ok := gcontext.GetRequestScopedFields(ctx)
	assert.Assert(t, ok, ok)
	assert.Assert(t, rsFields.TraceID == "", rsFields)
	assert.Assert(t, rsFields.RequestID != "", rsFields)
	assert.Assert(t, rsFields.CorrelationID != "", rsFields)
	assert.Assert(t, rsFields.CustomerAggregateID == "", rsFields)
	assert.Assert(t, rsFields.UserAggregateID == "", rsFields)
}