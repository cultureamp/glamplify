package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ShouldLogLevel(t *testing.T) {

	leveller := NewLevelMap()

	ok := leveller.ShouldLogLevel(DebugLevel, DebugLevel)
	assert.True(t, ok)
	ok = leveller.ShouldLogLevel(InfoLevel, DebugLevel)
	assert.False(t, ok)
	ok = leveller.ShouldLogLevel(DebugLevel, InfoLevel)
	assert.True(t, ok)
}

func Test_ShouldLogSeverity(t *testing.T) {

	leveller := NewLevelMap()

	ok := leveller.ShouldLogSeverity(DebugSev, DebugSev)
	assert.True(t, ok)
	ok = leveller.ShouldLogSeverity(InfoSev, DebugSev)
	assert.False(t, ok)
	ok = leveller.ShouldLogSeverity(DebugSev, InfoSev)
	assert.True(t, ok)
}

func Test_StringToLevel(t *testing.T) {

	leveller := NewLevelMap()

	level := leveller.StringToLevel(DebugSev)
	assert.Equal(t, DebugLevel, level)
	level = leveller.StringToLevel(InfoSev)
	assert.Equal(t, InfoLevel, level)
	level = leveller.StringToLevel(WarnSev)
	assert.Equal(t, WarnLevel, level)
	level = leveller.StringToLevel(ErrorSev)
	assert.Equal(t, ErrorLevel, level)
	level = leveller.StringToLevel(FatalSev)
	assert.Equal(t, FatalLevel, level)
	level = leveller.StringToLevel(AuditSev)
	assert.Equal(t, AuditLevel, level)
	level = leveller.StringToLevel("bad")
	assert.Equal(t, DebugLevel, level)
}
