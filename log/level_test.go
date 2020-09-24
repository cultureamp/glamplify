package log

import (
	"gotest.tools/assert"
	"testing"
)

func Test_ShouldLogLevel(t *testing.T) {

	leveller := NewLevelMap()

	ok := leveller.ShouldLogLevel(DebugLevel, DebugLevel)
	assert.Assert(t, ok, ok)
	ok = leveller.ShouldLogLevel(InfoLevel, DebugLevel)
	assert.Assert(t, !ok, ok)
	ok = leveller.ShouldLogLevel(DebugLevel, InfoLevel)
	assert.Assert(t, ok, ok)
}

func Test_ShouldLogSeverity(t *testing.T) {

	leveller := NewLevelMap()

	ok := leveller.ShouldLogSeverity(DebugSev, DebugSev)
	assert.Assert(t, ok, ok)
	ok = leveller.ShouldLogSeverity(InfoSev, DebugSev)
	assert.Assert(t, !ok, ok)
	ok = leveller.ShouldLogSeverity(DebugSev, InfoSev)
	assert.Assert(t, ok, ok)
}

func Test_StringToLevel(t *testing.T) {

	leveller := NewLevelMap()

	level := leveller.StringToLevel(DebugSev)
	assert.Assert(t, level == DebugLevel)
	level = leveller.StringToLevel(InfoSev)
	assert.Assert(t, level == InfoLevel)
	level = leveller.StringToLevel(WarnSev)
	assert.Assert(t, level == WarnLevel)
	level = leveller.StringToLevel(ErrorSev)
	assert.Assert(t, level == ErrorLevel)
	level = leveller.StringToLevel(FatalSev)
	assert.Assert(t, level == FatalLevel)
	level = leveller.StringToLevel("bad")
	assert.Assert(t, level == DebugLevel)

}