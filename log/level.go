package log

// Leveller manages the conversion of log levels to and from string to int
type Leveller struct {
	stol map[string]int
}

const (
	DebugLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	AuditLevel
)

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

func (sev Leveller) StringToLevel(severity string) int {
	level, ok := sev.stol[severity]
	if ok {
		return level
	}

	return DebugLevel
}

func (sev Leveller) ShouldLogSeverity(level string, severity string) bool {
	l := sev.StringToLevel(level)
	s := sev.StringToLevel(severity)

	return sev.ShouldLogLevel(l, s)
}

func (sev Leveller) ShouldLogLevel(level int, severity int) bool {
	if severity >= level {
		return true
	}
	return false
}