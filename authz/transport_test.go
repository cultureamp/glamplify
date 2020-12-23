package authz

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Http_Transport_New(t *testing.T) {
	client := NewHTTPTransport()
	assert.NotNil(t, client)
}

func Test_Http_Transport_New_With_Config(t *testing.T) {
	client := NewHTTPTransport(func(config *HTTPConfig) {
		assert.Equal(t, 10 * time.Second, config.ClientTimeout)
		assert.Equal(t, 5 * time.Second, config.DialerTimeout)
		assert.Equal(t, 5 * time.Second, config.TLSHandshakeTimeout)
	})
	assert.NotNil(t, client)
}

func Test_Http_Transport_Post_Error(t *testing.T){
	client := NewHTTPTransport(func(config *HTTPConfig) {
		assert.Equal(t, 10 * time.Second, config.ClientTimeout)
		assert.Equal(t, 5 * time.Second, config.DialerTimeout)
		assert.Equal(t, 5 * time.Second, config.TLSHandshakeTimeout)
	})

	response, err := client.Post(nil, "http://error.local", "application/json", bytes.NewBuffer([]byte("{}")))
	assert.Error(t, err)
	assert.EqualError(t, err, "net/http: nil Context")
	assert.Nil(t, response)

	ctx := context.Background()
	response, err = client.Post(ctx, "http://error.local", "application/json", bytes.NewBuffer([]byte("{}")))
	assert.Error(t, err)
	assert.Nil(t, response)

	// This gives difference answers on Windows vs Linux
	// assert.EqualError(t, err, "Post \"http://error.local\": dial tcp: lookup error.local: no such host")
	// So just go with this for now...
	assert.Contains(t, err.Error(), "dial tcp")
}