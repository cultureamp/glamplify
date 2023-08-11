package env

const (

	// *** AuthZ Environment Variables ***
	// AuthzClientTimeoutEnv = "AUTHZ_CLIENT_TIMEOUT_IN_MS"
	AuthzClientTimeoutEnv = "AUTHZ_CLIENT_TIMEOUT_IN_MS"
	// AuthzCacheDurationEnv = "AUTHZ_CACHE_DURATION_IN_SEC"
	AuthzCacheDurationEnv = "AUTHZ_CACHE_DURATION_IN_SEC"
	// AuthzDialerTimeoutEnv = "AUTHZ_DIALER_TIMEOUT_IN_MS"
	AuthzDialerTimeoutEnv = "AUTHZ_DIALER_TIMEOUT_IN_MS"
	// AuthzTLSTimeoutEnv    = "AUTHZ_TLS_TIMEOUT_IN_MS"
	AuthzTLSTimeoutEnv = "AUTHZ_TLS_TIMEOUT_IN_MS"

	// *** AWS Environment Variables ***
	// AwsProfileEnv = "AWS_PROFILE"
	AwsProfileEnv = "AWS_PROFILE"
	// AwsRegionEnv = "AWS_REGION"
	AwsRegionEnv = "AWS_REGION"
	// AwsAccountIDEnv  = "AWS_ACCOUNT_ID"
	AwsAccountIDEnv = "AWS_ACCOUNT_ID"
	// AwsXrayEnv = "XRAY_LOGGING"
	AwsXrayEnv = "XRAY_LOGGING"

	// *** Cache Environment Variables ***
	// CacheDurationEnv = "CACHE_DURATION_IN_SEC"
	CacheDurationEnv = "CACHE_DURATION_IN_SEC"

	// *** Datadog Environment Variables ***
	// DatadogAPIEnvVar        = "DD_API_KEY"
	DatadogAPIEnvVar = "DD_API_KEY"
	// DatadogLogEndpoint   = "DD_LOG_ENDPOINT"
	DatadogLogEndpoint = "DD_LOG_ENDPOINT"
	// DatadogEnv           = "DD_ENV"
	DatadogEnv = "DD_ENV"
	// DatadogService       = "DD_SERVICE"
	DatadogService = "DD_SERVICE"
	// DatadogVersion       = "DD_VERSION"
	DatadogVersion = "DD_VERSION"
	// DatadogAgentHost     = "DD_AGENT_HOST"
	DatadogAgentHost = "DD_AGENT_HOST"
	// DatadogStatsdPort = "DD_DOGSTATSD_PORT"
	DatadogStatsdPort = "DD_DOGSTATSD_PORT"
	// DatadogTimeout       = "DD_TIMEOUT"
	DatadogTimeout = "DD_TIMEOUT"
	// DatadogSite          = "DD_SITE"
	DatadogSite = "DD_SITE"
	// DatadogLogLevel      = "DD_LOG_LEVEL"
	DatadogLogLevel = "DD_LOG_LEVEL"

	// *** Log Environment Variables ***
	// Level            = "LOG_LEVEL"
	LogLevel = "LOG_LEVEL"
	// OmitEmpty        = "LOG_OMITEMPTY"
	LogOmitEmpty = "LOG_OMITEMPTY"
	// UseColours       = "LOG_COLOURS"
	LogUseColours = "LOG_COLOURS"

	// *** Sentry Environment Variables ***
	// SentryDsnEnv  = "SENTRY_DSN"
	SentryDsnEnv = "SENTRY_DSN"
	// SentryFlushTimeoutInMsEnv = "SENTRY_FLUSH_TIMEOUT_IN_MS"
	SentryFlushTimeoutInMsEnv = "SENTRY_FLUSH_TIMEOUT_IN_MS"

	// *** Global Environment Variables ***
	// AppNameEnv = "APP"
	AppNameEnv = "APP"
	// AppVerEnv = "APP_VERSION"
	AppVerEnv = "APP_VERSION"
	// AppEnv = "APP_ENV"
	AppEnv = "APP_ENV"
	// AppFarmEnv = "FARM"
	AppFarmEnv = "FARM"
	// ProductEnv = "PRODUCT"
	ProductEnv = "PRODUCT"
	// AppFarmLegacyEnv = "APP_ENV"
	AppFarmLegacyEnv = "APP_ENV"
)
