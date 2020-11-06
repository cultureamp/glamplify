package log_test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cultureamp/glamplify/context"
	"gotest.tools/assert"
	"strings"
	"testing"

	"github.com/cultureamp/glamplify/log"
)

func Test_Segment_Debug(t *testing.T) {

	memBuffer, logger := getTestLogger()

	properties := log.Fields{
		"string": "hello world",
		"int":    123,
	}
	logger.Event("something_happened").Fields(properties).Debug("not sure what is going on!")

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "something_happened")
	assertContainsString(t, msg, "severity", "DEBUG")
	assertContainsString(t, msg, "message", "not sure what is going on!")
	assertContainsString(t, msg, "string", "hello world")
	assertContainsInt(t, msg, "int", 123)
}

func Test_Segment_Info(t *testing.T) {

	memBuffer, logger := getTestLogger()

	properties := log.Fields{
		"string": "hello world",
		"int":    123,
	}
	logger.Event("something_happened").Fields(properties).Info("not sure what is going on!")

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "something_happened")
	assertContainsString(t, msg, "severity", "INFO")
	assertContainsString(t, msg, "message", "not sure what is going on!")
	assertContainsString(t, msg, "string", "hello world")
	assertContainsInt(t, msg, "int", 123)
}

func Test_Segment_Audit(t *testing.T) {

	memBuffer, logger := getTestLogger()

	properties := log.Fields{
		"string": "hello world",
		"int":    123,
	}
	logger.Event("something_happened").Fields(properties).Audit("not sure what is going on!")

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "something_happened")
	assertContainsString(t, msg, "severity", "AUDIT")
	assertContainsString(t, msg, "message", "not sure what is going on!")
	assertContainsString(t, msg, "string", "hello world")
	assertContainsInt(t, msg, "int", 123)
}

func Test_Segment_Warn(t *testing.T) {

	memBuffer, logger := getTestLogger()

	properties := log.Fields{
		"string": "hello world",
		"int":    123,
	}
	logger.Event("something_happened").Fields(properties).Warn("not sure what is going on!")

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "something_happened")
	assertContainsString(t, msg, "severity", "WARN")
	assertContainsString(t, msg, "message", "not sure what is going on!")
	assertContainsString(t, msg, "string", "hello world")
	assertContainsInt(t, msg, "int", 123)
}

func Test_Segment_Error(t *testing.T) {

	memBuffer, logger := getTestLogger()

	properties := log.Fields{
		"string": "hello world",
		"int":    123,
	}
	logger.Event("something_happened").Fields(properties).Error(errors.New("not sure what is going on"))

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "something_happened")
	assertContainsString(t, msg, "severity", "ERROR")
	assertContainsString(t, msg, "error", "not sure what is going on")
	assertContainsString(t, msg, "string", "hello world")
	assertContainsInt(t, msg, "int", 123)
}

func Test_Segment_Fatal(t *testing.T) {

	memBuffer, logger := getTestLogger()

	properties := log.Fields{
		"string": "hello world",
		"int":    123,
	}

	defer func() {
		if r := recover(); r != nil {
			msg := memBuffer.String()
			assertContainsString(t, msg, "event", "something_happened")
			assertContainsString(t, msg, "severity", "FATAL")
			assertContainsString(t, msg, "error", "not sure what is going on")
			assertContainsString(t, msg, "string", "hello world")
			assertContainsInt(t, msg, "int", 123)
		}
	}()

	logger.Event("something_happened").Fields(properties).Fatal(errors.New("not sure what is going on"))
}

func Test_Segment_WithNoFields(t *testing.T) {

	memBuffer, logger := getTestLogger()

	logger.Event("something_happened").Info("nothing to write home about")

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "something_happened")
	assertContainsString(t, msg, "severity", "INFO")
	assertContainsString(t, msg, "message", "nothing to write home about")
}

func Test_Segment_WithMultipleFields(t *testing.T) {

	memBuffer, logger := getTestLogger()

	logger.Event("something_happened").Fields(log.Fields{
		"string": "hello world",
	}).Fields(log.Fields{
		"int":    123,
	}).Info("nothing to write home about")

	msg := memBuffer.String()
	assertContainsString(t, msg, "event", "something_happened")
	assertContainsString(t, msg, "severity", "INFO")
	assertContainsString(t, msg, "message", "nothing to write home about")
	assertContainsString(t, msg, "string", "hello world")
	assertContainsInt(t, msg, "int", 123)
}


func getTestLogger() (*bytes.Buffer, *log.Logger) {
	rsFields := context.RequestScopedFields{
		TraceID:             "1-2-3",
		RequestID:           "4-5-6",
		CustomerAggregateID: "abc",
		UserAggregateID:     "xyz",
	}

	memBuffer := &bytes.Buffer{}
	writer := log.NewWriter(func(conf *log.WriterConfig) {
		conf.Output = memBuffer
	})
	logger := log.NewWitCustomWriter(rsFields, writer)
	return memBuffer, logger
}

func assertContainsString(t *testing.T, log string, key string, val string) {
	// Check that the keys and values are in the log line
	find := fmt.Sprintf("\"%s\":\"%s\"", key, val)
	assert.Assert(t, strings.Contains(log, find), "Expected '%s' in '%s'", find, log)
}

func assertContainsInt(t *testing.T, log string, key string, val int) {
	// Check that the keys and values are in the log line
	find := fmt.Sprintf("\"%s\":%d", key, val)
	assert.Assert(t, strings.Contains(log, find), "Expected '%s' in '%s'", find, log)
}