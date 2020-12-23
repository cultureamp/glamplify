package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_SetGet(t *testing.T) {

	cache := New()
	assert.NotNil(t, cache)

	val := "a value"
	cache.Set("k1", val, 1*time.Second)

	x, found := cache.Get("k1")
	assert.True(t, found)
	v := x.(string)
	assert.Equal(t, "a value", v)
}

func Test_SetGet_Expiry(t *testing.T) {

	cache := New()
	assert.NotNil(t, cache)

	val := "a value"
	cache.Set("k1", val, 1*time.Second)

	time.Sleep(2*time.Second)

	_, found := cache.Get("k1")
	assert.False(t, found)
}
