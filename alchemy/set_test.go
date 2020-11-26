package alchemy

import (
	"math/rand"
	"testing"
	"time"

	"gotest.tools/assert"
)

func Test_New_BitSet(t *testing.T) {
	lhs := newBitSet(testCauldron)
	assert.Assert(t, lhs != nil, lhs)
}

func Test_BitSet_And(t *testing.T) {
	lhs := newBitSet(testCauldron)
	rhs := newBitSet(testCauldron)

	// And two empty sets
	result, err := lhs.And(rhs)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, result != nil, result)
	assert.Assert(t, result.Count() == 0, result.Count())

	// One empty, the other with values
	empty := newBitSet(testCauldron)
	lhs = newBitSet(testCauldron)
	rhs = newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := Long(rand.Int63n(TestSetMaxSize))
		rhs.SetBit(bitIdx)
		lhs.SetBit(bitIdx)
	}

	result, err = empty.And(rhs)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, result != nil, result)
	assert.Assert(t, result.Count() == 0, result.Count())

	result, err = rhs.And(empty)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, result != nil, result)
	assert.Assert(t, result.Count() == 0, result.Count())

	result, err = lhs.And(rhs)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, result != nil, result)
	assert.Assert(t, result.Count() == lhs.Count(), result.Count())
	assert.Assert(t, result.Count() == rhs.Count(), result.Count())
}

func Test_BitSet_AndCount(t *testing.T) {
	lhs := newBitSet(testCauldron)
	rhs := newBitSet(testCauldron)

	// And two empty sets
	count, err := lhs.AndCount(rhs)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, count == 0, count)

	// One empty, the other with values
	empty := newBitSet(testCauldron)
	lhs = newBitSet(testCauldron)
	rhs = newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := Long(rand.Int63n(TestSetMaxSize))
		rhs.SetBit(bitIdx)
		lhs.SetBit(bitIdx)
	}

	count, err = empty.AndCount(rhs)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, count == 0, count)

	count, err = rhs.AndCount(empty)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, count == 0, count)

	count, err = lhs.AndCount(rhs)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, count == lhs.Count(), count)
	assert.Assert(t, count == rhs.Count(), count)
}

func Test_BitSet_Or(t *testing.T) {
	lhs := newBitSet(testCauldron)
	rhs := newBitSet(testCauldron)

	// Or two empty sets
	result, err := lhs.Or(rhs)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, result != nil, result)
	assert.Assert(t, result.Count() == 0, result.Count())

	// One empty, the other with values
	empty := newBitSet(testCauldron)
	lhs = newBitSet(testCauldron)
	rhs = newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := Long(rand.Int63n(TestSetMaxSize))
		rhs.SetBit(bitIdx)
		lhs.SetBit(bitIdx)
	}

	result, err = empty.Or(rhs)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, result != nil, result)
	assert.Assert(t, result.Count() == rhs.Count(), result.Count())

	result, err = rhs.Or(empty)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, result != nil, result)
	assert.Assert(t, result.Count() == rhs.Count(), result.Count())

	result, err = lhs.Or(rhs)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, result != nil, result)
	assert.Assert(t, result.Count() == lhs.Count(), result.Count())
	assert.Assert(t, result.Count() == rhs.Count(), result.Count())
}

func Test_BitSet_OrCount(t *testing.T) {
	lhs := newBitSet(testCauldron)
	rhs := newBitSet(testCauldron)

	// Or two empty sets
	count, err := lhs.OrCount(rhs)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, count == 0, count)

	// One empty, the other with values
	empty := newBitSet(testCauldron)
	lhs = newBitSet(testCauldron)
	rhs = newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := Long(rand.Int63n(TestSetMaxSize))
		rhs.SetBit(bitIdx)
		lhs.SetBit(bitIdx)
	}

	count, err = empty.OrCount(rhs)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, count == rhs.Count(), count)

	count, err = rhs.OrCount(empty)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, count == rhs.Count(), count)

	count, err = lhs.OrCount(rhs)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, count == lhs.Count(), count)
	assert.Assert(t, count == rhs.Count(), count)
}

func Test_BitSet_Not(t *testing.T) {
	set := newBitSet(testCauldron)

	// Not empty set
	result, err := set.Not()
	assert.Assert(t, err == nil, err)
	assert.Assert(t, result != nil, result)
	assert.Assert(t, result.Count() == TestSetMaxSize, result.Count())

	// Add some values
	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	set = newBitSet(testCauldron)
	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := Long(rand.Int63n(TestSetMaxSize))
		set.SetBit(bitIdx)
	}

	countBeforeNot := set.Count()

	result, err = set.Not()
	assert.Assert(t, err == nil, err)
	assert.Assert(t, result != nil, result)
	assert.Assert(t, result.Count() == (TestSetMaxSize-countBeforeNot), result.Count())
}

func Test_BitSet_NotCount(t *testing.T) {
	set := newBitSet(testCauldron)

	// Not empty set
	count, err := set.NotCount()
	assert.Assert(t, err == nil, err)
	assert.Assert(t, count == TestSetMaxSize, count)

	// Add some values
	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	set = newBitSet(testCauldron)
	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := Long(rand.Int63n(TestSetMaxSize))
		set.SetBit(bitIdx)
	}

	countBeforeNot := set.Count()

	count, err = set.NotCount()
	assert.Assert(t, err == nil, err)
	assert.Assert(t, count == (TestSetMaxSize-countBeforeNot), count)
}