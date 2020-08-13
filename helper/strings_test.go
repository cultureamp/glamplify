package helper

import (
	"gotest.tools/assert"
	"testing"
)

func Test_ToSnakeCase(t *testing.T) {

	sc := ToSnakeCase("hello")
	assert.Assert(t,  sc == "hello", "was: '%s'", sc)

	sc = ToSnakeCase("requestID")
	assert.Assert(t,  sc == "request_id", "was: '%s'", sc)

	sc = ToSnakeCase("LastGC")
	assert.Assert(t,  sc == "last_gc", "was: '%s'", sc)

	sc = ToSnakeCase("request_id")
	assert.Assert(t,  sc == "request_id", "was: '%s'", sc)

	sc = ToSnakeCase("something happened")
	assert.Assert(t,  sc == "something_happened", "was: '%s'", sc)

	sc = ToSnakeCase(" with  added  spaces")
	assert.Assert(t,  sc == "with_added_spaces", "was: '%s'", sc)

	sc = ToSnakeCase(" And  WITH  Capitals  ")
	assert.Assert(t,  sc == "and_with_capitals","was: '%s'", sc)
}

func Test_Redact(t *testing.T) {
	s := ""
	r := Redact(s)
	assert.Assert(t, r == "******", r)

	s = "1234"
	r = Redact(s)
	assert.Assert(t, r == "******", r)

	s = "123456"
	r = Redact(s)
	assert.Assert(t, r == "******", r)

	s = "1234567"
	r = Redact(s)
	assert.Assert(t, r == "******7", r)

	s = "12345678"
	r = Redact(s)
	assert.Assert(t, r == "******78", r)

	s = "123456789"
	r = Redact(s)
	assert.Assert(t, r == "******789", r)

	s = "1234567890"
	r = Redact(s)
	assert.Assert(t, r == "******7890", r)

	s = "12345678901"
	r = Redact(s)
	assert.Assert(t, r == "******8901", r)

	s = "123456789012"
	r = Redact(s)
	assert.Assert(t, r == "******9012", r)

	s = "1234567890123"
	r = Redact(s)
	assert.Assert(t, r == "******0123", r)

}

func Benchmark_ToSnakeCase(b *testing.B) {

	sa := []string {"hello", "requestID", "request_id", "something happened"}

	for n := 0; n < b.N; n++ {
		for _, s := range sa {
			ToSnakeCase(s)
		}
	}
}
