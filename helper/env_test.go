package helper

import (
	"gotest.tools/assert"
	"os"
	"testing"
)

func Test_GetEnvString(t *testing.T) {
	os.Setenv("TEST_STRING", "string")
	defer os.Unsetenv("TEST_STRING")

	val := GetEnvString("should_not_exist_env_var", "fallback")
	assert.Assert(t, val == "fallback", val)

	val = GetEnvString("TEST_STRING", "fallback")
	assert.Assert(t, val == "string", val)
}

func Test_GetEnvInt(t *testing.T) {
	os.Setenv("TEST_INT", "123")
	defer os.Unsetenv("TEST_INT")

	val := GetEnvInt("should_not_exist_env_var", 42)
	assert.Assert(t, val == 42, val)

	val = GetEnvInt("TEST_INT", 6)
	assert.Assert(t, val == 123, val)
}

func Test_GetEnvBool(t *testing.T) {
	os.Setenv("TEST_BOOL", "true")
	defer os.Unsetenv("TEST_BOOL")

	val := GetEnvBool("should_not_exist_env_var", false)
	assert.Assert(t, val == false, val)

	val = GetEnvBool("TEST_BOOL", false)
	assert.Assert(t, val == true, val)
}
