package opa

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"gotest.tools/assert"
)

func Test_OPAClient_New(t *testing.T) {
	client := NewOPAClient("dummy", mockHttpClient{})
	assert.Assert(t, client != nil, client)
}

func Test_OPAClient_Throw_Error(t *testing.T) {
	ctx := context.Background()
	client := NewOPAClient("dummy", mockHttpClient{throwError: true})

	response, err := client.EvaluateBooleanPolicy(ctx, "test.policy.name", IdentityRequest{}, InputRequest{})
	assert.Assert(t, err != nil, err)
	assert.Assert(t, response == nil, response)
}

func Test_OPAClient_Return_Empty(t *testing.T) {
	ctx := context.Background()
	client := NewOPAClient("dummy", mockHttpClient{returnEmpty: true})

	response, err := client.EvaluateBooleanPolicy(ctx, "test.policy.name", IdentityRequest{}, InputRequest{})
	assert.Assert(t, err != nil, err)
	assert.Assert(t, response == nil, response)
}

func Test_OPAClient_Return_Bad_JSON(t *testing.T) {
	ctx := context.Background()
	client := NewOPAClient("dummy", mockHttpClient{returnBadJson: true})

	response, err := client.EvaluateBooleanPolicy(ctx, "test.policy.name", IdentityRequest{}, InputRequest{})
	assert.Assert(t, err != nil, err)
	assert.Assert(t, response == nil, response)
}

func Test_OPAClient_Return_Not_Allowed(t *testing.T) {
	ctx := context.Background()
	client := NewOPAClient("dummy", mockHttpClient{returnNotAllowed: true})

	response, err := client.EvaluateBooleanPolicy(ctx, "test.policy.name", IdentityRequest{}, InputRequest{})
	assert.Assert(t, err == nil, err)
	assert.Assert(t, response != nil, response)
	assert.Assert(t, response.Allow == false, response)
}

func Test_OPAClient_Return_Allowed(t *testing.T) {
	ctx := context.Background()
	client := NewOPAClient("dummy", mockHttpClient{})

	response, err := client.EvaluateBooleanPolicy(ctx, "test.policy.name", IdentityRequest{}, InputRequest{})
	assert.Assert(t, err == nil, err)
	assert.Assert(t, response != nil, response)
	assert.Assert(t, response.Allow == true, response)
}

type mockHttpClient struct {
	throwError    bool
	returnEmpty   bool
	returnBadJson bool
	returnNotAllowed bool
}

func (client mockHttpClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	if client.throwError {
		return nil, errors.New("internal server error")
	}

	if client.returnEmpty {
		response := &http.Response{
			Body:       http.NoBody,
			StatusCode: http.StatusInternalServerError,
			Status:     http.StatusText(http.StatusInternalServerError),
		}
		return response, nil
	}

	if client.returnBadJson {
		response := &http.Response{
			Body:       ioutil.NopCloser(strings.NewReader("{hello world}")), // r type is io.ReadCloser,
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
		}
		return response, nil
	}

	if client.returnNotAllowed {
		postBody := `{
  		"result": [{
			"policy": "test.policy.name",
    		"allow": false,
    		"status": "ok"
  		}]}`

		response := &http.Response{
			Body:       ioutil.NopCloser(strings.NewReader(postBody)), // r type is io.ReadCloser,
			StatusCode: http.StatusOK,
			Status:     http.StatusText(http.StatusOK),
			Header:     http.Header{"Cache-Control": []string{"no-cache"}},
		}
		return response, nil
	}

	postBody := `{
  		"result": [{
			"policy": "test.policy.name",
    		"allow": true,
    		"status": "ok"
  		}]}`

	response := &http.Response{
		Body:       ioutil.NopCloser(strings.NewReader(postBody)), // r type is io.ReadCloser,
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Header:     http.Header{"Cache-Control": []string{"max-age=60"}},
	}
	return response, nil
}
