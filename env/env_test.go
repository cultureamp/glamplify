package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetEnvString(t *testing.T) {
	os.Setenv("TEST_STRING", "string")
	defer os.Unsetenv("TEST_STRING")

	val := GetString("should_not_exist_env_var", "fallback")
	assert.Equal(t, "fallback", val)

	val = GetString("TEST_STRING", "fallback")
	assert.Equal(t, "string", val)
}

func Test_GetEnvInt(t *testing.T) {
	os.Setenv("TEST_INT", "123")
	defer os.Unsetenv("TEST_INT")

	val := GetInt("should_not_exist_env_var", 42)
	assert.Equal(t, 42, val)

	val = GetInt("TEST_INT", 6)
	assert.Equal(t, 123, val)
}

func Test_GetEnvBool(t *testing.T) {
	os.Setenv("TEST_BOOL", "true")
	defer os.Unsetenv("TEST_BOOL")

	val := GetBool("should_not_exist_env_var", false)
	assert.False(t, val)

	val = GetBool("TEST_BOOL", false)
	assert.True(t, val)
}
