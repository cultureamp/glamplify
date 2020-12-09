package context_test

import (
	"github.com/cultureamp/glamplify/context"
	"github.com/cultureamp/glamplify/jwt"
	"gotest.tools/assert"
	"net/http"
	"testing"
)


func Test_RequestScope_AddGet(t *testing.T) {


	req, _ := http.NewRequest("GET", "*", nil)
	req = context.AddRequestScopedFieldsRequest(req, context.RequestScopedFields{
		TraceID:             "1-2-3",
		RequestID:           "7-8-9",
		CorrelationID:       "1-5-9",
		UserAggregateID:     "a-b-c",
		CustomerAggregateID: "xyz",
	})

	rsFields, ok := context.GetRequestScopedFieldsFromRequest(req)
	assert.Assert(t, ok, ok)
	assert.Assert(t, rsFields.TraceID == "1-2-3", rsFields)
	assert.Assert(t, rsFields.RequestID == "7-8-9", rsFields)
	assert.Assert(t, rsFields.CorrelationID == "1-5-9", rsFields)
	assert.Assert(t, rsFields.UserAggregateID == "a-b-c", rsFields)
	assert.Assert(t, rsFields.CustomerAggregateID == "xyz", rsFields)
}

func Test_Request_Wrap(t *testing.T) {

	req, err := http.NewRequest("GET", "*", nil)
	assert.Assert(t, err == nil, err)
	req.Header.Set(context.TraceIDHeader, "a-b-c")
	req.Header.Set(context.RequestIDHeader, "1-2-3")
	req.Header.Set(context.CorrelationIDHeader, "5-6-7")

	req, err = context.WrapRequest(req)
	assert.Assert(t, err != nil, err) // no pem key

	req, err = context.WrapRequestWithDecoder(req, nil) // no jwtDecoder
	assert.Assert(t, err != nil, err) // no pem key, but still populates the ID fields (but not the User/Account fields)
	rsFields, ok := context.GetRequestScopedFieldsFromRequest(req)
	assert.Assert(t, ok, ok)
	assert.Assert(t, rsFields.TraceID == "a-b-c", rsFields)
	assert.Assert(t, rsFields.RequestID == "1-2-3", rsFields)
	assert.Assert(t, rsFields.CorrelationID == "5-6-7", rsFields)
	assert.Assert(t, rsFields.UserAggregateID == "", rsFields)
	assert.Assert(t, rsFields.CustomerAggregateID == "", rsFields)
}

func Test_Request_WrapWithDecoder(t *testing.T) {
	jwt, err := jwt.NewDecoderFromPath("../jwt/jwt.rs256.key.development.pub")
	assert.Assert(t, err == nil, err)

	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiZXhwIjoxOTAzOTMwNzA0LCJpYXQiOjE1ODg1NzA3MDR9.XGm34FDIgtBFvx5yC2HTUu-cf3DaQI4TmIBVLx0H7y89oNVNWJaKA3dLvWS0oOZoYIuGhj6GzPREBEmou2f9JsUerqnc-_Tf8oekFZWU7kEfzu9ECBiSWPk7ljPJeZLbau62sSqD7rYb-m3v1mohqz4tKJ_7leWu9L1uHHliC7YGlSRl1ptVDllJjKXKjOg9ifeGSXDEMeU35KgCFwIwKdu8WmCTd8ztLSKEnLT1OSaRZ7MSpmHQ4wUZtS6qvhLBiquvHub9KdQmc4mYWLmfKdDiR5DH-aswJFGLVu3yisFRY8uSfeTPQRhQXd_UfdgifCTXdWTnCvNZT-BxULYG-5mlvAFu-JInTga_9-r-wHRzFD1SrcKjuECF7vUG8czxGNE4sPjFrGVyBxE6fzzcFsdrhdqS-LB_shVoG940fD-ecAhXQZ9VKgr-rmCvmxuv5vYI2HoMfg9j_-zeXkucKxvPYvDQZYMdeW4wFsUORliGplThoHEeRQxTX8d_gvZFCy_gGg0H57FmJwCRymWk9v29s6uyHUMor_r-e7e6ZlShFBrCPAghXL04S9IFJUxUv30wNie8aaSyvPuiTqCgGiEwF_20ZaHCgYX0zupdGm4pHTyJrx2wv31yZ4VZYt8tKjEW6-BlB0nxzLGk5OUN83vq-RzH-92WmY5kMndF6Jo"

	req, err := http.NewRequest("GET", "*", nil)
	assert.Assert(t, err == nil, err)
	req.Header.Set("Authorization", "Bearer " + token)
	req.Header.Set(context.TraceIDHeader, "a-b-c")
	req.Header.Set(context.RequestIDHeader, "1-2-3")
	req.Header.Set(context.CorrelationIDHeader, "5-6-7")

	req2, err := context.WrapRequestWithDecoder(req, jwt)
	assert.Assert(t, err == nil, err)
	rsFields, ok := context.GetRequestScopedFieldsFromRequest(req2)
	assert.Assert(t, ok, ok)
	assert.Assert(t, rsFields.TraceID == "a-b-c", rsFields)
	assert.Assert(t, rsFields.RequestID == "1-2-3", rsFields)
	assert.Assert(t, rsFields.CorrelationID == "5-6-7", rsFields)
	assert.Assert(t, rsFields.UserAggregateID != "", rsFields)
	assert.Assert(t, rsFields.CustomerAggregateID != "", rsFields)
}
