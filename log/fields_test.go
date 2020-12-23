package log_test

import (
	"testing"
	"time"

	"github.com/cultureamp/glamplify/log"
	"github.com/stretchr/testify/assert"
)

func TestFields_Success(t *testing.T) {
	entries := log.Fields{
		"aString": "hello world",
		"aInt":    123,
	}
	assert.NotNil(t, entries)

	ok, err := entries.ValidateNewRelic()
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestFields_Merge_Duration(t *testing.T) {
	d := time.Millisecond * 456
	durations := log.NewDurationFields(d)

	tt := durations["time_taken"]
	assert.Equal(t, "P0.456S", tt)
	ttms := durations["time_taken_ms"]
	assert.Equal(t, int64(456), ttms)

	entries := log.Fields{
		"aString": "hello world",
		"aInt":    123,
	}
	entries = entries.Merge(durations)

	tt = entries["time_taken"]
	assert.Equal(t, "P0.456S", tt)
	ttms = entries["time_taken_ms"]
	assert.Equal(t, int64(456), ttms)
}

func TestFields_InvalidType_Failed(t *testing.T) {
	dict := map[string]int{
		"key1": 1,
	}
	entries := log.Fields{
		"aMap": dict,
	}
	assert.NotNil(t, entries)

	ok, err := entries.ValidateNewRelic()
	assert.NotNil(t, err)
	assert.False(t, ok)
}

func TestFields_NilValue_Failed(t *testing.T) {
	dict := map[string]interface{}{
		"key1": nil,
	}
	entries := log.Fields{
		"aMap": dict,
		"akey": nil,
	}
	assert.NotNil(t, entries)

	ok, err := entries.ValidateNewRelic()
	assert.NotNil(t, err)
	assert.False(t, ok)
}

func TestFields_StringToLong_Failed(t *testing.T) {
	entries := log.Fields{
		"aString": "big_long_string_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890_1234567890",
	}
	assert.NotNil(t, entries)

	ok, err := entries.ValidateNewRelic()
	assert.NotNil(t, err)
	assert.False(t, ok)
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

	str := fields.ToJSON(false)
	assert.Equal(t, "{\"key_string\":\"abc\"}", str)
}

func TestFields_ToTags(t *testing.T) {
	fields := log.Fields{
		"key_string": "abc",
		"key_int":    1,
		"key_float":  3.14,
		"key_field": log.Fields{
			"sub_key_string": "xyz",
			"sub_key_int":    5,
			"sub_key_float":  6.28,
		},
	}

	tags := fields.ToTags(false)
	assert.Equal(t, 6, len(tags))
	assert.Contains(t, tags, "key_string:abc")
	assert.Contains(t, tags, "key_int:1")
	assert.Contains(t, tags, "key_float:3.14")
	assert.Contains(t, tags, "sub_key_string:xyz")
	assert.Contains(t, tags, "sub_key_int:5")
	assert.Contains(t, tags, "sub_key_float:6.28")
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
		fields.ToSnakeCase().ToJSON(false)
	}
}
