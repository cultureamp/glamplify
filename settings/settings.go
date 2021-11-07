package settings

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/cultureamp/glamplify/env"
)

// set as ldflags by build
var (
	appName     = "glamplify"
	appVersion  = "0.0.0-local"
	buildNumber = "0"
	branch      = "local"
	commit      = "unknown-local"
)

const (
	upstreamTimeoutMsEnv = "UPSTREAM_TIMEOUT_MS"
	jwtPublicKeyEnv      = "AUTH_PUBLIC_KEY"
	datadogEnabledEnv    = "DD_FLUSH_TO_LOG"
)

// Settings describes the application configuration
type Settings struct {
	env.Settings // anonymous field, inherit all the base env.Settings fields (eg. App, AppEnv, Farm, etc) and functions (eg. IsProduction(), IsRunningLocal(), etc)

	BuildNumber string
	Branch      string
	Commit      string

	MongoDatabase string
	MongoHost     string
	MongoUsername string
	MongoPassword string

	UpstreamTimeout time.Duration

	JwtPublicKey   string
	DatadogEnabled bool
}

// SetupEnvironment adds some environment variables that the logging package looks for.
func SetupEnvironment() {
	mustSet(os.Setenv(env.AppNameEnv, appName))
	mustSet(os.Setenv(env.AppVerEnv, appVersion))
}

// NewSettings gathers configuration from the environment as required.
func NewSettings() (*Settings, error) {
	gSettings := env.NewSettings()
	gSettings.App = appName
	gSettings.AppVersion = appVersion

	upstreamTimeoutMs := env.GetInt(upstreamTimeoutMsEnv, 300)
	if upstreamTimeoutMs < 1 {
		return nil, errors.New("environment variable " + upstreamTimeoutMsEnv + " is less than 1")
	}

	jwtPublicKey := os.Getenv(jwtPublicKeyEnv)
	if jwtPublicKey == "" {
		path := ".development-keys/web-gateway/jwt/public.pem"
		key, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("environment variable %s was empty, tried reading development key %s: %w", jwtPublicKeyEnv, path, err)
		}
		jwtPublicKey = string(key)
	}

	settings := &Settings{
		Settings:        *gSettings,
		BuildNumber:     buildNumber,
		Branch:          branch,
		Commit:          commit,
		UpstreamTimeout: time.Duration(upstreamTimeoutMs) * time.Millisecond,
		JwtPublicKey:    jwtPublicKey,
		DatadogEnabled:  env.GetBool(datadogEnabledEnv, false),
	}

	return settings, nil
}

func mustSet(err error) {
	if err != nil {
		panic(err)
	}
}
