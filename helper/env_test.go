package helper

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetEnvString(t *testing.T) {
	os.Setenv("TEST_STRING", "string")
	defer os.Unsetenv("TEST_STRING")

	val := GetEnvString("should_not_exist_env_var", "fallback")
	assert.Equal(t, "fallback", val)

	val = GetEnvString("TEST_STRING", "fallback")
	assert.Equal(t, "string", val)
}

func Test_GetEnvInt(t *testing.T) {
	os.Setenv("TEST_INT", "123")
	defer os.Unsetenv("TEST_INT")

	val := GetEnvInt("should_not_exist_env_var", 42)
	assert.Equal(t, 42, val)

	val = GetEnvInt("TEST_INT", 6)
	assert.Equal(t, 123, val)
}

func Test_GetEnvBool(t *testing.T) {
	os.Setenv("TEST_BOOL", "true")
	defer os.Unsetenv("TEST_BOOL")

	val := GetEnvBool("should_not_exist_env_var", false)
	assert.False(t, val)

	val = GetEnvBool("TEST_BOOL", false)
	assert.True(t, val)
}
