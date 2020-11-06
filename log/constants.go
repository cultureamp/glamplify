package log

const (
	Unknown = "unknown"
	Empty   = ""

	// RFC3339Milli is the standard RFC3339 format with added milliseconds
	RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

	// JSON LOG KEYS
	// List of standard keys used for logging as per https://cultureamp.atlassian.net/wiki/spaces/TV/pages/959939199/Logging
	TraceID             = "trace_id"
	RequestID           = "request_id"
	CorrelationID       = "correlation_id"
	Time                = "time"
	Event               = "event"
	Product             = "product"
	App                 = "app"
	Farm                = "farm"
	AppVer              = "app_version"
	Severity            = "severity"
	AwsRegion           = "aws_region"
	AwsAccountID        = "aws_account_id"
	Resource            = "resource"
	Os                  = "os"
	Customer            = "customer"
	User                = "user"
	Exception           = "exception"
	Message             = "message"
	Properties          = "properties"
	TimeTaken           = "time_taken"
	TimeTakenMS         = "time_taken_ms"
	MemoryUsed          = "memory_used"
	MemoryAvail         = "memory_available"
	ItemsProcessed      = "items_processed"
	TotalItemsProcessed = "total_items_processed"
	TotalItemsRequested = "total_items_requested"

	// Severity Values
	DebugSev = "DEBUG"
	InfoSev  = "INFO"
	WarnSev  = "WARN"
	ErrorSev = "ERROR"
	FatalSev = "FATAL"
	AuditSev = "AUDIT"

	// ENVIRONMENT VARIABLES
	Level            = "LOG_LEVEL"
	OmitEmpty        = "LOG_OMITEMPTY"
	UseColours       = "LOG_COLOURS"
	ProductEnv       = "PRODUCT"
	AppNameEnv       = "APP"
	AppFarmLegacyEnv = "APP_ENV"
	AppFarmEnv       = "FARM"
	AppVerEnv        = "APP_VERSION"
	AwsRegionEnv     = "AWS_REGION"
	AwsAccountIDEnv  = "AWS_ACCOUNT_ID"
)
