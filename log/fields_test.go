package log_test

import (
	"github.com/cultureamp/glamplify/log"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestFields_Success(t *testing.T) {
	entries := log.Fields{
		"aString": "hello world",
		"aInt":    123,
	}
	assert.Assert(t, entries != nil, entries)

	ok, err := entries.ValidateNewRelic()
	assert.Assert(t, ok, ok)
	assert.Assert(t, err == nil, err)
}

func TestFields_Merge_Duration(t *testing.T) {
	d := time.Millisecond * 456
	durations := log.NewDurationFields(d)

	tt := durations["time_taken"]
	assert.Assert(t, tt == "P0.456S", tt)
	ttms := durations["time_taken_ms"]
	assert.Assert(t, ttms == int64(456), ttms)

	entries := log.Fields{
		"aString": "hello world",
		"aInt":    123,
	}
	entries = entries.Merge(durations)

	tt = entries["time_taken"]
	assert.Assert(t, tt == "P0.456S", tt)
	ttms = entries["time_taken_ms"]
	assert.Assert(t, ttms == int64(456), ttms)
}

func TestFields_InvalidType_Failed(t *testing.T) {
	dict := map[string]int{
		"key1": 1,
	}
	entries := log.Fields{
		"aMap": dict,
	}
	assert.Assert(t, entries != nil, entries)

	ok, err := entries.ValidateNewRelic()
	assert.Assert(t, !ok, ok)
	assert.Assert(t, err != nil, err)
}

func TestFields_NilValue_Failed(t *testing.T) {
	dict := map[string]interface{}{
		"key1": nil,
	}
	entries := log.Fields{
		"aMap": dict,
		"akey": nil,
	}
	assert.Assert(t, entries != nil, entries)

	ok, err := entries.ValidateNewRelic()
	assert.Assert(t, !ok, ok)
	assert.Assert(t, err != nil, err)
}

func TestFields_StringToLong_Failed(t *testing.T) {
	entries := log.Fields{
		"aString": "big_long_string_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890",
	}
	assert.Assert(t, entries != nil, entries)

	ok, err := entries.ValidateNewRelic()
	assert.Assert(t, !ok, ok)
	assert.Assert(t, err != nil, err)
}

func TestFields_InvalidValues_ToJSON(t *testing.T) {
	fields := log.Fields{
		"key_string": "abc",
		"key_func": func() int64 {
			var l int64 = 123
			return l
		},
		"key_chan": make(chan string),
	}

	str := fields.ToJson(false)
	assert.Assert(t, str == "{\"key_string\":\"abc\"}", str)
}

func TestFields_ToTags(t *testing.T) {
	fields := log.Fields{
		"key_string": "abc",
		"key_int": 1,
		"key_float": 3.14,
		"key_field": log.Fields{
			"sub_key_string": "xyz",
			"sub_key_int": 5,
			"sub_key_float": 6.28,
		},
	}

	tags := fields.ToTags(false)
	assert.Assert(t, len(tags) == 6, tags)
	assert.Assert(t, tags[0] == "key_string:abc", tags)
	assert.Assert(t, tags[1] == "key_int:1", tags)
	assert.Assert(t, tags[2] == "key_float:3.14", tags)
	assert.Assert(t, tags[3] == "sub_key_string:xyz", tags)
	assert.Assert(t, tags[4] == "sub_key_int:5", tags)
	assert.Assert(t, tags[5] == "sub_key_float:6.28", tags)
}

func Benchmark_FieldsToJSON(b *testing.B) {

	fields := log.Fields{
		"string":        "hello",
		"int":           123,
		"float":         42.48,
		"string2":       "hello world",
		"string3 space": "world",
	}

	for n := 0; n < b.N; n++ {
		fields.ToSnakeCase().ToJson(false)
	}
}
