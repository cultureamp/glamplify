package alchemy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New_Stack(t *testing.T) {
	s := newLinkedListStack()
	assert.NotNil(t, s)
}

func Test_LinkedListStack(t *testing.T) {
	s := newLinkedListStack()
	assert.True(t, s.isEmpty())

	s.push(1)
	assert.False(t, s.isEmpty())

	id, err := s.pop()
	assert.Nil(t, err)
	assert.Equal(t, 1, id)
	assert.True(t, s.isEmpty())

	_, err = s.pop()
	assert.NotNil(t, err)

	// push and pop a number of elements
	for i := 0; i < TestNumberOfBits; i++ {
		s.push(uint64(i))
	}
	for i := TestNumberOfBits - 1; i >= 0; i-- {
		id, err := s.pop()
		assert.Nil(t, err)
		assert.Equal(t, uint64(i), id)
	}

	_, err = s.pop()
	assert.NotNil(t, err)
}

func Benchmark_LinkedListStack(b *testing.B) {
	s := newLinkedListStack()

	for n := 0; n < b.N; n++ {
		for i := 0; i < TestNumberOfBits; i++ {
			s.push(uint64(n))
		}
		for i := 0; i < TestNumberOfBits; i++ {
			s.pop()
		}
	}
}


