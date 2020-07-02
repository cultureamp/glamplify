package log

import (
	"os"
)

type systemLogLevel struct {
	sysLogLevel int
	stol        map[string]int
}

const (
	DebugLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)


func newSystemLogLevel() *systemLogLevel {

	table := map[string]int{
		DebugSev: DebugLevel,
		InfoSev:  InfoLevel,
		WarnSev:  WarnLevel,
		ErrorSev: ErrorLevel,
		FatalSev: FatalLevel,
	}

	level, ok := os.LookupEnv(Level)
	if !ok {
		level = DebugSev
	}
	logLevel, found := table[level]
	if !found {
		logLevel = DebugLevel
	}

	return &systemLogLevel{
		sysLogLevel: logLevel,
		stol:        table,
	}
}

func (sev systemLogLevel) stringToLevel(severity string) int {
	level, ok := sev.stol[severity]
	if ok {
		return level
	}

	return DebugLevel
}

func (sev systemLogLevel) shouldLog(severity int) bool {

	if severity >= sev.sysLogLevel {
		return true
	}

	return false
}

