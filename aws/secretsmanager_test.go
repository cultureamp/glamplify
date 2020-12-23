package aws

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/stretchr/testify/assert"
)

func Test_GetSecretParam_MissingKey_NoCache(t *testing.T) {

	sm := NewSecretsManager(func(config *SecretsManagerConfig) {
		config.Profile = "default"
		config.CacheErrorsAsEmpty = false
	})
	assert.NotNil(t, sm)

	// Missing Key
	val, err := sm.Get("/this/should/not/exist/secret_key")
	assert.NotNil(t, err)
	assert.Empty(t, val)

	_, ok := err.(awserr.Error)
	assert.True(t, ok)
}

func Test_GetSecretParam_MissingKey_WithCache(t *testing.T) {

	sm := NewSecretsManager(func(config *SecretsManagerConfig) {
		config.Profile = "default"
		config.CacheErrorsAsEmpty = true
		config.CacheDuration = 1 * time.Minute
	})
	assert.NotNil(t, sm)

	// Missing Key
	val, err := sm.Get("/this/should/not/exist/secret_key")
	assert.NotNil(t, err)
	assert.Empty(t, val)

	val, err = sm.Get("/this/should/not/exist/secret_key")
	assert.Nil(t, err)
	assert.Empty(t, val)
}



// TODO - what is a good key & env to use for unit tests?

