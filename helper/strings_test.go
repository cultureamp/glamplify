package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ToSnakeCase(t *testing.T) {

	sc := ToSnakeCase("hello")
	assert.Equal(t, "hello", sc)

	sc = ToSnakeCase("requestID")
	assert.Equal(t, "request_id", sc)

	sc = ToSnakeCase("LastGC")
	assert.Equal(t, "last_gc", sc)

	sc = ToSnakeCase("request_id")
	assert.Equal(t, "request_id", sc)

	sc = ToSnakeCase("something happened")
	assert.Equal(t, "something_happened", sc)

	sc = ToSnakeCase(" with  added  spaces")
	assert.Equal(t, "with_added_spaces", sc)

	sc = ToSnakeCase(" And  WITH  Capitals  ")
	assert.Equal(t, "and_with_capitals", sc)
}

func Test_Redact(t *testing.T) {
	s := ""
	r := Redact(s)
	assert.Equal(t, "******", r)

	s = "1234"
	r = Redact(s)
	assert.Equal(t, "******", r)

	s = "123456"
	r = Redact(s)
	assert.Equal(t, "******", r)

	s = "1234567"
	r = Redact(s)
	assert.Equal(t, "******7", r)

	s = "12345678"
	r = Redact(s)
	assert.Equal(t, "******78", r)

	s = "123456789"
	r = Redact(s)
	assert.Equal(t, "******789", r)

	s = "1234567890"
	r = Redact(s)
	assert.Equal(t, "******7890", r)

	s = "12345678901"
	r = Redact(s)
	assert.Equal(t, "******8901", r)

	s = "123456789012"
	r = Redact(s)
	assert.Equal(t, "******9012", r)

	s = "1234567890123"
	r = Redact(s)
	assert.Equal(t, "******0123", r)

}

func Benchmark_ToSnakeCase(b *testing.B) {

	sa := []string{"hello", "requestID", "request_id", "something happened"}

	for n := 0; n < b.N; n++ {
		for _, s := range sa {
			ToSnakeCase(s)
		}
	}
}
