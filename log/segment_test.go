package log_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/cultureamp/glamplify/context"
	"github.com/cultureamp/glamplify/log"
	"github.com/stretchr/testify/assert"
)

func Test_Segment_Debug(t *testing.T) {

	memBuffer, logger := getTestLogger()

	properties := log.Fields{
		"string": "hello world",
		"int":    123,
	}
	logger.Event("something_happened").Fields(properties).Debug("not sure what is going on!")

	msg := memBuffer.String()
	assert.Contains(t, msg, "\"event\":\"something_happened\"")
	assert.Contains(t, msg, "\"severity\":\"DEBUG\"")
	assert.Contains(t, msg, "\"message\":\"not sure what is going on!\"")
	assert.Contains(t, msg, "\"string\":\"hello world\"")
	assert.Contains(t, msg, "\"int\":123")
}

func Test_Segment_Info(t *testing.T) {

	memBuffer, logger := getTestLogger()

	properties := log.Fields{
		"string": "hello world",
		"int":    123,
	}
	logger.Event("something_happened").Fields(properties).Info("not sure what is going on!")

	msg := memBuffer.String()
	assert.Contains(t, msg, "\"event\":\"something_happened\"")
	assert.Contains(t, msg, "\"severity\":\"INFO\"")
	assert.Contains(t, msg, "\"message\":\"not sure what is going on!\"")
	assert.Contains(t, msg, "\"string\":\"hello world\"")
	assert.Contains(t, msg, "\"int\":123")
}


func Test_Segment_Warn(t *testing.T) {

	memBuffer, logger := getTestLogger()

	properties := log.Fields{
		"string": "hello world",
		"int":    123,
	}
	logger.Event("something_happened").Fields(properties).Warn("not sure what is going on!")

	msg := memBuffer.String()
	assert.Contains(t, msg, "\"event\":\"something_happened\"")
	assert.Contains(t, msg, "\"severity\":\"WARN\"")
	assert.Contains(t, msg, "\"message\":\"not sure what is going on!\"")
	assert.Contains(t, msg, "\"string\":\"hello world\"")
	assert.Contains(t, msg, "\"int\":123")
}

func Test_Segment_Error(t *testing.T) {

	memBuffer, logger := getTestLogger()

	properties := log.Fields{
		"string": "hello world",
		"int":    123,
	}
	logger.Event("something_happened").Fields(properties).Error(errors.New("not sure what is going on"))

	msg := memBuffer.String()
	assert.Contains(t, msg, "\"event\":\"something_happened\"")
	assert.Contains(t, msg, "\"severity\":\"ERROR\"")
	assert.Contains(t, msg, "\"error\":\"not sure what is going on\"")
	assert.Contains(t, msg, "\"string\":\"hello world\"")
	assert.Contains(t, msg, "\"int\":123")
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
			assert.Contains(t, msg, "\"event\":\"something_happened\"")
			assert.Contains(t, msg, "\"severity\":\"FATAL\"")
			assert.Contains(t, msg, "\"error\":\"not sure what is going on\"")
			assert.Contains(t, msg, "\"string\":\"hello world\"")
			assert.Contains(t, msg, "\"int\":123")
		}
	}()

	logger.Event("something_happened").Fields(properties).Fatal(errors.New("not sure what is going on"))
}

func Test_Segment_Audit(t *testing.T) {

	memBuffer, logger := getTestLogger()

	properties := log.Fields{
		"string": "hello world",
		"int":    123,
	}
	logger.Event("something_happened").Fields(properties).Audit("not sure what is going on!")

	msg := memBuffer.String()
	assert.Contains(t, msg, "\"event\":\"something_happened\"")
	assert.Contains(t, msg, "\"severity\":\"AUDIT\"")
	assert.Contains(t, msg, "\"message\":\"not sure what is going on!\"")
	assert.Contains(t, msg, "\"string\":\"hello world\"")
	assert.Contains(t, msg, "\"int\":123")
}

func Test_Segment_WithNoFields(t *testing.T) {

	memBuffer, logger := getTestLogger()

	logger.Event("something_happened").Info("nothing to write home about")

	msg := memBuffer.String()
	assert.Contains(t, msg, "\"event\":\"something_happened\"")
	assert.Contains(t, msg, "\"severity\":\"INFO\"")
	assert.Contains(t, msg, "\"message\":\"nothing to write home about\"")
}

func Test_Segment_WithMultipleFields(t *testing.T) {

	memBuffer, logger := getTestLogger()

	logger.Event("something_happened").Fields(log.Fields{
		"string": "hello world",
	}).Fields(log.Fields{
		"int":    123,
	}).Info("nothing to write home about")

	msg := memBuffer.String()
	assert.Contains(t, msg, "\"event\":\"something_happened\"")
	assert.Contains(t, msg, "\"severity\":\"INFO\"")
	assert.Contains(t, msg, "\"message\":\"nothing to write home about\"")
	assert.Contains(t, msg, "\"string\":\"hello world\"")
	assert.Contains(t, msg, "\"int\":123")
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