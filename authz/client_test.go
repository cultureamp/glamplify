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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_OPAClient_New(t *testing.T) {
	client := NewClient("dummy", mockHTTPClient{})
	assert.NotNil(t, client)
}

func Test_OPAClient_Throw_Error(t *testing.T) {
	ctx := context.Background()
	client := NewClient("dummy", mockHTTPClient{throwError: true},  func(config *Config) {
		config.Timeout = 100 * time.Millisecond
	})

	response, err := client.EvaluateBooleanPolicy(ctx, "test.policy.name", IdentityRequest{}, InputRequest{})
	assert.NotNil(t, err)
	assert.Nil(t, response)
}

func Test_OPAClient_Sleep(t *testing.T) {
	ctx := context.Background()
	client := NewClient("dummy", mockHTTPClient{sleep: true}, func(config *Config) {
		config.Timeout = 100 * time.Millisecond
	})

	response, err := client.EvaluateBooleanPolicy(ctx, "test.policy.name", IdentityRequest{}, InputRequest{})
	assert.NotNil(t, err)
	assert.Equal(t, "context deadline exceeded", err.Error())
	assert.Nil(t, response)
}

func Test_OPAClient_Return_Empty(t *testing.T) {
	ctx := context.Background()
	client := NewClient("dummy", mockHTTPClient{returnEmpty: true}, func(config *Config) {
		config.Timeout = 100 * time.Millisecond
	})

	response, err := client.EvaluateBooleanPolicy(ctx, "test.policy.name", IdentityRequest{}, InputRequest{})
	assert.NotNil(t, err)
	assert.Nil(t, response)
}

func Test_OPAClient_Return_Bad_JSON(t *testing.T) {
	ctx := context.Background()
	client := NewClient("dummy", mockHTTPClient{returnBadJSON: true},  func(config *Config) {
		config.Timeout = 100 * time.Millisecond
	})

	response, err := client.EvaluateBooleanPolicy(ctx, "test.policy.name", IdentityRequest{}, InputRequest{})
	assert.NotNil(t, err)
	assert.Nil(t, response)
}

func Test_OPAClient_Return_Not_Allowed(t *testing.T) {
	ctx := context.Background()
	client := NewClient("dummy", mockHTTPClient{returnNotAllowed: true},  func(config *Config) {
		config.Timeout = 100 * time.Millisecond
	})

	response, err := client.EvaluateBooleanPolicy(ctx, "test.policy.name", IdentityRequest{}, InputRequest{})
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.False(t, response.Allow)
}

func Test_OPAClient_Return_Allowed(t *testing.T) {
	ctx := context.Background()
	client := NewClient("dummy", mockHTTPClient{},  func(config *Config) {
		config.Timeout = 100 * time.Millisecond
	})

	response, err := client.EvaluateBooleanPolicy(ctx, "test.policy.name", IdentityRequest{}, InputRequest{})
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Allow)
}

type mockHTTPClient struct {
	mock.Mock
	throwError       bool
	sleep            bool
	returnEmpty      bool
	returnBadJSON    bool
	returnNotAllowed bool
}

func (client mockHTTPClient) Post(ctx context.Context, url string, contentType string, body io.Reader) (resp *http.Response, err error) {
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

	if client.returnBadJSON {
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
