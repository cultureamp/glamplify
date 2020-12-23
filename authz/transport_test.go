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
		assert.Equal(t, 5 * time.Second, config.DialerTimeout)
		assert.Equal(t, 5 * time.Second, config.TLSHandshakeTimeout)
	})
	assert.NotNil(t, client)
}