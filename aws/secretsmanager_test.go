package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/stretchr/testify/assert"
)

func Test_GetSecretParam_MissingKey(t *testing.T) {

	sm := NewSecretsManager("default")
	assert.NotNil(t, sm)

	// Missing Key
	val, err := sm.Get("/this/should/not/exist/secret_key")
	assert.NotNil(t, err)
	assert.Empty(t, val)

	aerr, ok := err.(awserr.Error)
	assert.True(t, ok)
	assert.NotEmpty(t, aerr.Message())
}

// TODO - what is a good key to use for unit tests?

