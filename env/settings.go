package env

import (
	"encoding/json"
	"github.com/cultureamp/glamplify/helper"
	"os"
	"strings"
)

// Settings that drive behavior.
// Consumers of this library should inherit this by embedding this struct
// in their own Settings
type Settings struct {
	// Common environment variable values used by at least 80% of apps
	App           string `json:"app"`
	AppVersion    string `json:"app_version"`
	AppEnv        string `json:"app_env"`
	Farm          string `json:"farm"`
	Product       string `json:"product"`
	AwsProfile    string `json:"aws_profile"`
	AwsRegion     string `json:"aws_region"`
	AwsAccountID  string `json:"aws_account_id"`
	XrayLogging   bool   `json:"xray_logging"`
	DDApiKey      string `json:"dd_api_key"`
	SentryDSN     string `json:"sentry_dsn"`
	SentryFlushMs int    `json:"sentry_flush_ms"`
}

func NewSettings() *Settings {

	settings := &Settings{}

	settings.App = GetString(AppNameEnv, "authz-api")
	settings.AppVersion = GetString(AppVerEnv, "1.0.0")
	settings.AppEnv = GetString(AppEnv, "development")
	settings.Farm = GetString(AppFarmEnv, "production")
	settings.Product = os.Getenv(ProductEnv)
	settings.AwsProfile = GetString(AwsProfileEnv, "default")
	settings.AwsRegion = GetString(AwsRegionEnv, "us-west-2")
	settings.AwsAccountID = os.Getenv(AwsAccountIDEnv)
	settings.XrayLogging = GetBool(AwsXrayEnv, true)
	settings.DDApiKey = os.Getenv(DatadogApiKey)
	settings.SentryDSN = os.Getenv(SentryDsnEnv)
	settings.SentryFlushMs = GetInt(SentryFlushTimeoutInMsEnv, 50)

	return settings
}

// RedactedSettings returns redacted Settings
func (s Settings) RedactedSettings() *Settings {
	return &Settings{
		App:           s.App,
		AppVersion:    s.AppVersion,
		AppEnv:        s.AppEnv,
		Farm:          s.Farm,
		Product:       s.Product,
		AwsProfile:    s.AwsProfile,
		AwsRegion:     s.AwsRegion,
		AwsAccountID:  s.AwsAccountID,
		XrayLogging:   s.XrayLogging,
		DDApiKey:      helper.Redact(s.DDApiKey),
		SentryDSN:     helper.Redact(s.SentryDSN),
		SentryFlushMs: s.SentryFlushMs,
	}
}

// IsProduction returns true if "APP_ENV" == "local"
func (s Settings) IsProduction() bool {
	return s.AppEnv == "production"
}

// IsRunningInAWS returns true if "FARM" != "local"
func (s Settings) IsRunningInAWS() bool {
	return !s.IsRunningLocal()
}

// IsRunningLocal returns true if FARM" == "local"
func (s Settings) IsRunningLocal() bool {
	return s.Farm == "local"
}

// ToJSON returns Settings as a JSON string
func (s Settings) ToJSON() string {
	b, err := json.Marshal(s)
	if err != nil {
		// https://stackoverflow.com/questions/33903552/what-input-will-cause-golangs-json-marshal-to-return-an-error#:~:text=From%20the%20docs%3A,result%20in%20an%20infinite%20recursion.
		// should not happen with a valid Settings
		panic(err)
	}
	return string(b)
}

// ToRedactedJSON returns Settings as a redacted JSON string
func (s Settings) ToRedactedJSON() string {
	rs := s.RedactedSettings()
	return rs.ToJSON()
}

// ToString returns Settings as a string (not redacted)
func (s Settings) ToString() string {
	data := s.ToJSON()
	return strings.Replace(data, "\"", "", -1)
}

// ToRedactedString returns Settings as a redacted string
func (s Settings) ToRedactedString() string {
	rs := s.RedactedSettings()
	return rs.ToString()
}
