package alchemy

import (
	"gotest.tools/assert"
	"testing"
)

func Test_New_Stack(t *testing.T) {
	s := newLinkedListStack()
	assert.Assert(t, s != nil, s)
}

func Test_LinkedListStack(t *testing.T) {
	s := newLinkedListStack()
	assert.Assert(t, s.isEmpty(), s.isEmpty())

	s.push(1)
	assert.Assert(t, !s.isEmpty(), s.isEmpty())

	id, err := s.pop()
	assert.Assert(t, err == nil, err)
	assert.Assert(t, id == 1, id)
	assert.Assert(t, s.isEmpty(), s.isEmpty())

	_, err = s.pop()
	assert.Assert(t, err != nil, err)

	// push and pop a number of elements
	for i := 0; i < TestNumberOfBits; i++ {
		s.push(uint64(i))
	}
	for i := TestNumberOfBits - 1; i >= 0; i-- {
		id, err := s.pop()
		assert.Assert(t, err == nil, err)
		assert.Assert(t, id == uint64(i), id)
	}

	_, err = s.pop()
	assert.Assert(t, err != nil, err)
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


