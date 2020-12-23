package log

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WriteFields(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
		conf.OmitEmpty = false
	})

	writer.WriteFields(DebugSev, Fields{
		"system":       "system_value",
		"system_empty": "",
	}, Fields{
		"properties":       "properties_value",
		"properties_empty": "",
	})

	msg := memBuffer.String()
	assert.Contains(t, msg, "\"system\":\"system_value\"")
	assert.Contains(t, msg, "\"system_empty\":\"\"")
	assert.Contains(t, msg, "\"properties\":\"properties_value\"")
	assert.Contains(t, msg, "\"properties_empty\":\"\"")
}

func Test_WriteFields_OmitEmpty(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
		conf.OmitEmpty = true
	})

	writer.WriteFields(DebugSev,
		Fields{
			"system":       "system_value",
			"system_empty": "",
		}, Fields{
			"properties":       "properties_value",
			"properties_empty": "",
		})

	msg := memBuffer.String()
	assert.Contains(t, msg, "\"system\":\"system_value\"")
	assert.NotContains(t, msg, "system_empty")
	assert.Contains(t, msg, "\"properties\":\"properties_value\"")
	assert.NotContains(t, msg, "properties_empty")
}

func Test_WriteFields_IsEnabled(t *testing.T) {
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Level = InfoSev
	})

	ok := writer.IsEnabled(DebugSev)
	assert.False(t, ok)
	ok = writer.IsEnabled(InfoSev)
	assert.True(t, ok)
	ok = writer.IsEnabled(WarnSev)
	assert.True(t, ok)
	ok = writer.IsEnabled(ErrorSev)
	assert.True(t, ok)
	ok = writer.IsEnabled(FatalSev)
	assert.True(t, ok)
	ok = writer.IsEnabled(AuditSev)
	assert.True(t, ok)
}
