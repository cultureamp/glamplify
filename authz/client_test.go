package authz

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"gotest.tools/assert"
)

func Test_OPAClient_New(t *testing.T) {
	client := NewClient("dummy", mockHttpClient{})
	assert.Assert(t, client != nil, client)
}

func Test_OPAClient_Throw_Error(t *testing.T) {
	ctx := context.Background()
	client := NewClient("dummy", mockHttpClient{throwError: true},  func(config *Config) {
		config.Timeout = 100 * time.Millisecond
	})

	response, err := client.EvaluateBooleanPolicy(ctx, "test.policy.name", IdentityRequest{}, InputRequest{})
	assert.Assert(t, err != nil, err)
	assert.Assert(t, response == nil, response)
}

func Test_OPAClient_Sleep(t *testing.T) {
	ctx := context.Background()
	client := NewClient("dummy", mockHttpClient{sleep: true}, func(config *Config) {
		config.Timeout = 100 * time.Millisecond
	})

	response, err := client.EvaluateBooleanPolicy(ctx, "test.policy.name", IdentityRequest{}, InputRequest{})
	assert.Assert(t, err != nil, err)
	assert.Assert(t, err.Error() == "context deadline exceeded")
	assert.Assert(t, response == nil, response)
}

func Test_OPAClient_Return_Empty(t *testing.T) {
	ctx := context.Background()
	client := NewClient("dummy", mockHttpClient{returnEmpty: true}, func(config *Config) {
		config.Timeout = 100 * time.Millisecond
	})

	response, err := client.EvaluateBooleanPolicy(ctx, "test.policy.name", IdentityRequest{}, InputRequest{})
	assert.Assert(t, err != nil, err)
	assert.Assert(t, response == nil, response)
}

func Test_OPAClient_Return_Bad_JSON(t *testing.T) {
	ctx := context.Background()
	client := NewClient("dummy", mockHttpClient{returnBadJson: true},  func(config *Config) {
		config.Timeout = 100 * time.Millisecond
	})

	response, err := client.EvaluateBooleanPolicy(ctx, "test.policy.name", IdentityRequest{}, InputRequest{})
	assert.Assert(t, err != nil, err)
	assert.Assert(t, response == nil, response)
}

func Test_OPAClient_Return_Not_Allowed(t *testing.T) {
	ctx := context.Background()
	client := NewClient("dummy", mockHttpClient{returnNotAllowed: true},  func(config *Config) {
		config.Timeout = 100 * time.Millisecond
	})

	response, err := client.EvaluateBooleanPolicy(ctx, "test.policy.name", IdentityRequest{}, InputRequest{})
	assert.Assert(t, err == nil, err)
	assert.Assert(t, response != nil, response)
	assert.Assert(t, response.Allow == false, response)
}

func Test_OPAClient_Return_Allowed(t *testing.T) {
	ctx := context.Background()
	client := NewClient("dummy", mockHttpClient{},  func(config *Config) {
		config.Timeout = 100 * time.Millisecond
	})

	response, err := client.EvaluateBooleanPolicy(ctx, "test.policy.name", IdentityRequest{}, InputRequest{})
	assert.Assert(t, err == nil, err)
	assert.Assert(t, response != nil, response)
	assert.Assert(t, response.Allow == true, response)
}

type mockHttpClient struct {
	throwError       bool
	sleep            bool
	returnEmpty      bool
	returnBadJson    bool
	returnNotAllowed bool
}

func (client mockHttpClient) Post(ctx context.Context, url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	if client.throwError {
		return nil, errors.New("internal server error")
	}

	if client.sleep {
		time.Sleep(500 * time.Millisecond)
		return nil, ctx.Err()
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
