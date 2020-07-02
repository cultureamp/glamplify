package log

import (
	"gotest.tools/assert"
	"os"
	"testing"
)

func Test_Sev_Log(t *testing.T) {

	sev := newSystemLogLevel()

	ok := sev.shouldLog(DebugLevel)
	assert.Assert(t, ok, ok)
	ok = sev.shouldLog(InfoLevel)
	assert.Assert(t, ok, ok)
	ok = sev.shouldLog(WarnLevel)
	assert.Assert(t, ok, ok)
	ok = sev.shouldLog(ErrorLevel)
	assert.Assert(t, ok, ok)
	ok = sev.shouldLog(FatalLevel)
	assert.Assert(t, ok, ok)
}

func Test_Sev_Log_Env(t *testing.T) {

	os.Setenv("LOG_LEVEL", WarnSev)
	defer os.Unsetenv("LOG_LEVEL")

	sev := newSystemLogLevel()

	ok := sev.shouldLog(DebugLevel)
	assert.Assert(t, !ok, ok)
	ok = sev.shouldLog(InfoLevel)
	assert.Assert(t, !ok, ok)
	ok = sev.shouldLog(WarnLevel)
	assert.Assert(t, ok, ok)
	ok = sev.shouldLog(ErrorLevel)
	assert.Assert(t, ok, ok)
	ok = sev.shouldLog(FatalLevel)
	assert.Assert(t, ok, ok)
}

func Test_Sev_Log_Unknown(t *testing.T) {

	sev := newSystemLogLevel()

	ok := sev.shouldLog(-1)
	assert.Assert(t, !ok, ok)
}

func Test_Sev_Log_Env_Unknown(t *testing.T) {

	os.Setenv("LOG_LEVEL", "unknown")
	defer os.Unsetenv("LOG_LEVEL")

	sev := newSystemLogLevel()

	ok := sev.shouldLog(DebugLevel)
	assert.Assert(t, ok, ok)
	ok = sev.shouldLog(InfoLevel)
	assert.Assert(t, ok, ok)
	ok = sev.shouldLog(WarnLevel)
	assert.Assert(t, ok, ok)
	ok = sev.shouldLog(ErrorLevel)
	assert.Assert(t, ok, ok)
	ok = sev.shouldLog(FatalLevel)
	assert.Assert(t, ok, ok)
}