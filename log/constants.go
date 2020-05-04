package log

type EventCtxKey int

const (
	Unknown = "unknown"
	Empty   = ""

	// RFC3339Milli is the standard RFC3339 format with added milliseconds
	RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

	// JSON LOG KEYS
	// List of standard keys used for logging as per https://cultureamp.atlassian.net/wiki/spaces/TV/pages/959939199/Logging
	TraceID  = "trace_id"
	Time     = "time"
	Event    = "event"
	Product  = "product"
	App      = "app"
	AppVer   = "app_version"
	Severity = "severity"
	Region   = "region"
	Resource = "resource"
	Os       = "os"
	Customer = "customer"
	User                = "user"
	Exception           = "exception"
	Message             = "message"
	TimeTaken           = "time_taken"
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

	// ENVIRONMENT VARIABLES
	// List of  Environment Variables keys
	ProductEnv = "PRODUCT"
	AppEnv     = "APP"
	AppVerEnv  = "APP_VERSION"
	RegionEnv  = "REGION"

	// CONTEXT KEYS
	TraceIDCtx  EventCtxKey = iota
	CustomerCtx EventCtxKey = iota
	UserCtx     EventCtxKey = iota
)
