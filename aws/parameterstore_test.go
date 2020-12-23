package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetParam_MissingKey(t *testing.T) {

	ps := NewParameterStore("default")
	assert.NotNil(t, ps)

	// Missing Key
	val, err := ps.Get("/this/should/not/exist/secret_key")
	assert.Empty(t, val)
	assert.NotNil(t, err)

	// aerr, ok := err.(awserr.Error)
}

// TODO - what is a good key to use for unit tests?
