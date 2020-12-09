package authz

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func Test_HttpClient_New(t *testing.T) {
	client := NewHTTPTransport()
	assert.Assert(t, client != nil, client)
}

func Test_HttpClient_New_With_Config(t *testing.T) {
	client := NewHTTPTransport(func(config *HTTPConfig) {
		assert.Assert(t, config.ClientTimeout == 10 * time.Second, config.ClientTimeout)
		config.ClientTimeout = 1 * time.Second
		assert.Assert(t, config.DialerTimeout == 5 * time.Second, config.DialerTimeout)
		config.DialerTimeout = 1 * time.Second
		assert.Assert(t, config.TLSHandshakeTimeout == 5 * time.Second, config.TLSHandshakeTimeout)
		config.TLSHandshakeTimeout = 1 * time.Second
	})
	assert.Assert(t, client != nil, client)
}