package authz

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_HttpClient_New(t *testing.T) {
	client := NewHTTPTransport()
	assert.NotNil(t, client)
}

func Test_HttpClient_New_With_Config(t *testing.T) {
	client := NewHTTPTransport(func(config *HTTPConfig) {
		assert.Equal(t, 10 * time.Second, config.ClientTimeout)
		config.ClientTimeout = 1 * time.Second
		assert.Equal(t, 5 * time.Second, config.DialerTimeout)
		config.DialerTimeout = 1 * time.Second
		assert.Equal(t, 5 * time.Second, config.TLSHandshakeTimeout)
		config.TLSHandshakeTimeout = 1 * time.Second
	})
	assert.NotNil(t, client)
}