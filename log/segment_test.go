package log

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Segment_Debug(t *testing.T) {

	logger := New(rsFields)

	properties := Fields{
		"string": "hello world",
		"int":    123,
	}
	json := logger.Event("something_happened").Fields(properties).Debug("not sure what is going on!")

	assert.Contains(t, json, "\"event\":\"something_happened\"")
	assert.Contains(t, json, "\"severity\":\"DEBUG\"")
	assert.Contains(t, json, "\"message\":\"not sure what is going on!\"")
	assert.Contains(t, json, "\"string\":\"hello world\"")
	assert.Contains(t, json, "\"int\":123")
}

func Test_Segment_Info(t *testing.T) {

	logger := New(rsFields)

	properties := Fields{
		"string": "hello world",
		"int":    123,
	}
	json := logger.Event("something_happened").Fields(properties).Info("not sure what is going on!")

	assert.Contains(t, json, "\"event\":\"something_happened\"")
	assert.Contains(t, json, "\"severity\":\"INFO\"")
	assert.Contains(t, json, "\"message\":\"not sure what is going on!\"")
	assert.Contains(t, json, "\"string\":\"hello world\"")
	assert.Contains(t, json, "\"int\":123")
}

func Test_Segment_Warn(t *testing.T) {

	logger := New(rsFields)

	properties := Fields{
		"string": "hello world",
		"int":    123,
	}
	json := logger.Event("something_happened").Fields(properties).Warn("not sure what is going on!")

	assert.Contains(t, json, "\"event\":\"something_happened\"")
	assert.Contains(t, json, "\"severity\":\"WARN\"")
	assert.Contains(t, json, "\"message\":\"not sure what is going on!\"")
	assert.Contains(t, json, "\"string\":\"hello world\"")
	assert.Contains(t, json, "\"int\":123")
}

func Test_Segment_Error(t *testing.T) {

	logger := New(rsFields)

	properties := Fields{
		"string": "hello world",
		"int":    123,
	}
	json := logger.Event("something_happened").Fields(properties).Error(errors.New("not sure what is going on"))

	assert.Contains(t, json, "\"event\":\"something_happened\"")
	assert.Contains(t, json, "\"severity\":\"ERROR\"")
	assert.Contains(t, json, "\"error\":\"not sure what is going on\"")
	assert.Contains(t, json, "\"string\":\"hello world\"")
	assert.Contains(t, json, "\"int\":123")
}

func Test_Segment_Fatal(t *testing.T) {

	memBuffer := &bytes.Buffer{}
	writer := NewWriter(func(conf *WriterConfig) {
		conf.Output = memBuffer
	})
	logger := NewWitCustomWriter(rsFields, writer)

	properties := Fields{
		"string": "hello world",
		"int":    123,
	}

	defer func() {
		if r := recover(); r != nil {
			json := memBuffer.String()
			assert.Contains(t, json, "\"event\":\"something_happened\"")
			assert.Contains(t, json, "\"severity\":\"FATAL\"")
			assert.Contains(t, json, "\"error\":\"not sure what is going on\"")
			assert.Contains(t, json, "\"string\":\"hello world\"")
			assert.Contains(t, json, "\"int\":123")
		}
	}()

	logger.Event("something_happened").Fields(properties).Fatal(errors.New("not sure what is going on"))
}

func Test_Segment_Audit(t *testing.T) {

	logger := New(rsFields)

	properties := Fields{
		"string": "hello world",
		"int":    123,
	}
	json := logger.Event("something_happened").Fields(properties).Audit("not sure what is going on!")

	assert.Contains(t, json, "\"event\":\"something_happened\"")
	assert.Contains(t, json, "\"severity\":\"AUDIT\"")
	assert.Contains(t, json, "\"message\":\"not sure what is going on!\"")
	assert.Contains(t, json, "\"string\":\"hello world\"")
	assert.Contains(t, json, "\"int\":123")
}

func Test_Segment_WithNoFields(t *testing.T) {

	logger := New(rsFields)

	json := logger.Event("something_happened").Info("nothing to write home about")

	assert.Contains(t, json, "\"event\":\"something_happened\"")
	assert.Contains(t, json, "\"severity\":\"INFO\"")
	assert.Contains(t, json, "\"message\":\"nothing to write home about\"")
}

func Test_Segment_WithMultipleFields(t *testing.T) {

	logger := New(rsFields)

	json := logger.Event("something_happened").Fields(Fields{
		"string": "hello world",
	}).Fields(Fields{
		"int":    123,
	}).Info("nothing to write home about")

	assert.Contains(t, json, "\"event\":\"something_happened\"")
	assert.Contains(t, json, "\"severity\":\"INFO\"")
	assert.Contains(t, json, "\"message\":\"nothing to write home about\"")
	assert.Contains(t, json, "\"string\":\"hello world\"")
	assert.Contains(t, json, "\"int\":123")
}