package log

import "os"

type systemLogLevel struct {
	sysLogLevel int
	lookup      map[string]int
}

func newSystemLogLevel() *systemLogLevel {

	table := map[string]int {
		DebugSev: 0,
		InfoSev: 1,
		WarnSev: 2,
		ErrorSev: 3,
		FatalSev: 4,
	}

	level, ok := os.LookupEnv(Level)
	if !ok {
		level = DebugSev
	}
	logLevel, found := table[level]
	if !found {
		logLevel = 0
	}

	return &systemLogLevel{
		sysLogLevel: logLevel,
		lookup: table,
	}
}

func (sev systemLogLevel) shouldLog(severity string) bool {

	level, ok := sev.lookup[severity]
	if !ok {
		return false
	}

	if level >= sev.sysLogLevel {
		return true
	}

	return false
}
