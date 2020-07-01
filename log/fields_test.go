package log_test

import (
	"github.com/cultureamp/glamplify/log"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestEntries_Success(t *testing.T) {
	entries := log.Fields{
		"aString": "hello world",
		"aInt":    123,
	}
	assert.Assert(t, entries != nil, entries)

	ok, err := entries.ValidateNewRelic()
	assert.Assert(t, ok, ok)
	assert.Assert(t, err == nil, err)
}

func TestEntries_Merge_Duration(t *testing.T) {
	d := time.Millisecond * 456
	durations := log.NewDurationFields(d)

	tt :=  durations["time_taken"]
	assert.Assert(t, tt == "P0.456S", tt)
	ttms := durations["time_taken_ms"]
	assert.Assert(t, ttms == int64(456), ttms)

	entries := log.Fields{
		"aString": "hello world",
		"aInt":    123,
	}
	entries = entries.Merge(durations)

	tt =  entries["time_taken"]
	assert.Assert(t, tt == "P0.456S", tt)
	ttms =  entries["time_taken_ms"]
	assert.Assert(t, ttms == int64(456), ttms)
}

func TestEntries_InvalidType_Failed(t *testing.T) {
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

func TestEntries_NilValue_Failed(t *testing.T) {
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

func TestEntries_StringToLong_Failed(t *testing.T) {
	entries := log.Fields{
		"aString": "big_long_string_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890",
	}
	assert.Assert(t, entries != nil, entries)

	ok, err := entries.ValidateNewRelic()
	assert.Assert(t, !ok, ok)
	assert.Assert(t, err != nil, err)
}

func TestEntries_InvalidValues_ToJSON(t*testing.T) {
	fields := log.Fields{
		"key_string": "abc",
		"key_func": func() int64 {
			var  l int64 = 123
			return l},
		"key_chan": make(chan string),
	}

	str := fields.ToJson(false)
	assert.Assert(t, str=="{\"key_string\":\"abc\"}", str)
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