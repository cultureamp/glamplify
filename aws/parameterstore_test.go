package aws

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/stretchr/testify/assert"
)

func Test_GetParam_MissingKey_NoCache(t *testing.T) {

	ps := NewParameterStore(func(config *ParameterStoreConfig) {
		config.Profile = "default"
		config.CacheErrorsAsEmpty = false
	})
	assert.NotNil(t, ps)

	// Missing Key
	val, err := ps.Get("/this/should/not/exist/secret_key")
	assert.NotNil(t, err)
	assert.Empty(t, val)

	_, ok := err.(awserr.Error)
	assert.True(t, ok)
}

func Test_GetParam_MissingKey_With_Cache(t *testing.T) {

	ps := NewParameterStore(func(config *ParameterStoreConfig) {
		config.Profile = "default"
		config.CacheErrorsAsEmpty = true
		config.CacheDuration = 1 * time.Minute
	})
	assert.NotNil(t, ps)

	// Missing Key
	val, err := ps.Get("/this/should/not/exist/secret_key")
	assert.NotNil(t, err)
	assert.Empty(t, val)

	// should be cached
	val, err = ps.Get("/this/should/not/exist/secret_key")
	assert.Nil(t, err)
	assert.Empty(t, val)
}

// TODO - what is a good key & env to use for unit tests?
