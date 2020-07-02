package log

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	gcontext "github.com/cultureamp/glamplify/context"
)

// SystemValues
type SystemValues struct {

}

func DurationAsISO8601(duration time.Duration) string {
	return fmt.Sprintf("P%gS", duration.Seconds())
}

func newSystemValues() *SystemValues {
	return &SystemValues{}
}

func (df SystemValues) getSystemValues(rsFields gcontext.RequestScopedFields, event string, severity string) Fields {
	fields := Fields{
		Time:     df.timeNow(RFC3339Milli),
		Event:    event,
		Resource: df.hostName(),
		Os:       df.targetOS(),
		Severity: severity,
	}

	fields = df.getMandatoryFields(rsFields, fields)
	fields = df.getEnvFields(fields)

	return fields
}

func (df SystemValues) getErrorValues(err error, fields Fields) Fields {
	errorMessage := strings.TrimSpace(err.Error())

	stats := &debug.GCStats{}
	buf := debug.Stack()
	debug.ReadGCStats(stats)

	fields[Exception] = Fields{
		"error":    errorMessage,
		"trace":    string(buf),
		"gc_stats": Fields{
			"last_gc": stats.LastGC,
			"num_gc": stats.NumGC,
			"pause_total": stats.PauseTotal,
			"pause_history": stats.Pause,
			"pause_end": stats.PauseEnd,
			"page_quantiles": stats.PauseQuantiles,
		},
	}

	return fields
}

func (df SystemValues) getEnvFields(fields Fields) Fields {

	fields = df.addEnvFieldIfMissing(Product, ProductEnv, fields)
	fields = df.addEnvFieldIfMissing(App, AppEnv, fields)
	fields = df.addEnvFieldIfMissing(AppVer, AppVerEnv, fields)
	fields = df.addEnvFieldIfMissing(AwsRegion, AwsRegionEnv, fields)
	fields = df.addEnvFieldIfMissing(AwsAccountID, AwsAccountIDEnv, fields)

	return fields
}

func (df SystemValues) getMandatoryFields(rsFields gcontext.RequestScopedFields, fields Fields) Fields {

	fields = df.addMandatoryFieldIfMissing(TraceID, rsFields.TraceID, fields)
	fields = df.addMandatoryFieldIfMissing(RequestID, rsFields.RequestID, fields)
	fields = df.addMandatoryFieldIfMissing(CorrelationID, rsFields.CorrelationID, fields)
	fields = df.addMandatoryFieldIfMissing(Customer, rsFields.CustomerAggregateID, fields)
	fields = df.addMandatoryFieldIfMissing(User, rsFields.UserAggregateID, fields)

	return fields
}

func (df SystemValues) addEnvFieldIfMissing(fieldName string, osVar string, fields Fields) Fields {

	// If it contains it already, all good!
	if _, ok := fields[fieldName]; ok {
		return fields
	}

	// otherwise get env value from OS
	prod := os.Getenv(osVar)
	fields[fieldName] = prod

	return fields
}

func (df SystemValues) addMandatoryFieldIfMissing(fieldName string, fieldValue string, fields Fields) Fields {

	// If it contains it already, all good!
	if _, ok := fields[fieldName]; ok {
		return fields
	}

	fields[fieldName] = fieldValue
	return fields
}

func (df SystemValues) timeNow(format string) string {
	return time.Now().UTC().Format(format)
}

var host string
var hostOnce sync.Once

func (df SystemValues) hostName() string {

	var err error
	hostOnce.Do(func() {
		host, err = os.Hostname()
		if err != nil {
			host = Unknown
		}
	})

	return host
}

func (df SystemValues) targetOS() string {
	return runtime.GOOS
}
