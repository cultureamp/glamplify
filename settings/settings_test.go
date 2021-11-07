package settings

import (
	"os"
	"testing"
	"time"

	"github.com/cultureamp/glamplify/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEndpointSettings(t *testing.T) {
	unsetEnvironmentVariables()
	defer unsetEnvironmentVariables()

	os.Setenv(upstreamTimeoutMsEnv, "123")
	os.Setenv(jwtPublicKeyEnv, "jwt-public-key")

	sut, err := NewSettings()

	require.Nil(t, err)
	require.NotNil(t, sut)

	assert.Equal(t, time.Millisecond*123, sut.UpstreamTimeout)
	assert.Equal(t, "jwt-public-key", sut.JwtPublicKey)
	assert.False(t, sut.DatadogEnabled)
}

func TestDatadogEnabled(t *testing.T) {
	unsetEnvironmentVariables()
	defer unsetEnvironmentVariables()

	os.Setenv(jwtPublicKeyEnv, "jwt-public-key")
	os.Setenv(datadogEnabledEnv, "true")

	sut, err := NewSettings()

	require.Nil(t, err)
	require.NotNil(t, sut)

	assert.True(t, sut.DatadogEnabled)
}

func TestInvalidTimeoutValue(t *testing.T) {
	unsetEnvironmentVariables()
	defer unsetEnvironmentVariables()

	os.Setenv(upstreamTimeoutMsEnv, "-123")

	_, err := NewSettings()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "environment variable UPSTREAM_TIMEOUT_MS is less than 1")
}

func TestAppVariables(t *testing.T) {
	defer func() {
		os.Unsetenv(env.AppNameEnv)
		os.Unsetenv(env.AppVerEnv)
	}()

	SetupEnvironment()

	assert.Equal(t, "glamplify", os.Getenv(env.AppNameEnv))
	assert.Equal(t, "0.0.0-local", os.Getenv(env.AppVerEnv))
}

func unsetEnvironmentVariables() {
	os.Unsetenv(upstreamTimeoutMsEnv)
	os.Unsetenv(jwtPublicKeyEnv)
}
