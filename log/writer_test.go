package log

import (
	"bytes"
	"fmt"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func Test_WriteFields(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})

	writer.WriteFields(DebugLevel,
		Fields{
		"system":       "system_value",
		"system_empty": "",
	}, Fields{
		"properties":       "properties_value",
		"properties_empty": "",
	})

	msg := memBuffer.String()
	assertStringContains(t, msg, "system", "system_value")
	assertStringContains(t, msg, "system_empty", "")
	assertStringContains(t, msg, "properties", "properties_value")
	assertStringContains(t, msg, "properties_empty", "")
}

func Test_WriteFields_Filtered(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
		conf.OmitEmpty = true
	})

	writer.WriteFields(DebugLevel,
		Fields{
		"system":       "system_value",
		"system_empty": "",
	}, Fields{
		"properties":       "properties_value",
		"properties_empty": "",
	})

	msg := memBuffer.String()
	assertStringContains(t, msg, "system", "system_value")
	assertKeyMissing(t, msg, "system_empty")
	assertStringContains(t, msg, "properties", "properties_value")
	assertKeyMissing(t, msg, "properties_empty")
}

func assertStringContains(t *testing.T, log string, key string, val string) {
	// Check that the keys and values are in the log line
	find := fmt.Sprintf("\"%s\":\"%s\"", key, val)
	assert.Assert(t, strings.Contains(log, find), "Expected '%s' in '%s'", find, log)
}

func assertKeyMissing(t *testing.T, log string, key string) {
	// Check that the keys and values are in the log line
	find := fmt.Sprintf("\"%s\"", key)
	assert.Assert(t, !strings.Contains(log, find), "Expected '%s' in '%s'", find, log)
}
