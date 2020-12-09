package log

// Leveller manages the conversion of log levels to and from string to int
type Leveller struct {
	stol map[string]int
}

const (
	// DebugLevel = iota
	DebugLevel = iota
	// InfoLevel = iota + 1
	InfoLevel
	// WarnLevel = iota + 2
	WarnLevel
	// ErrorLevel = iota + 3
	ErrorLevel
	// FatalLevel = iota + 4
	FatalLevel
	// AuditLevel = iota + 5
	AuditLevel
)

// NewLevelMap creates a Leveller map
func NewLevelMap() *Leveller {

	table := map[string]int{
		DebugSev: DebugLevel,
		InfoSev:  InfoLevel,
		WarnSev:  WarnLevel,
		ErrorSev: ErrorLevel,
		FatalSev: FatalLevel,
		AuditSev: AuditLevel,
	}

	return &Leveller{
		stol: table,
	}
}

// StringToLevel given a string severity returns the int value
func (sev Leveller) StringToLevel(severity string) int {
	level, ok := sev.stol[severity]
	if ok {
		return level
	}

	return DebugLevel
}

// ShouldLogSeverity given the current level and a severity returns true if should be logged, false otherwise
func (sev Leveller) ShouldLogSeverity(level string, severity string) bool {
	l := sev.StringToLevel(level)
	s := sev.StringToLevel(severity)

	return sev.ShouldLogLevel(l, s)
}

// ShouldLogLevel given the current level and a severity returns true if should be logged, false otherwise
func (sev Leveller) ShouldLogLevel(level int, severity int) bool {
	if severity >= level {
		return true
	}
	return false
}