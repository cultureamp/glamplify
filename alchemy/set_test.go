package alchemy

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_New_BitSet(t *testing.T) {
	lhs := newBitSet(testCauldron)
	assert.NotNil(t, lhs, lhs)
}

func Test_BitSet_And(t *testing.T) {
	lhs := newBitSet(testCauldron)
	rhs := newBitSet(testCauldron)

	// And two empty sets
	result, err := lhs.And(rhs)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint64(0), result.Count())

	// One empty, the other with values
	empty := newBitSet(testCauldron)
	lhs = newBitSet(testCauldron)
	rhs = newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := uint64(rand.Int63n(TestSetSize))
		rhs.SetBit(bitIdx)
		lhs.SetBit(bitIdx)
	}

	result, err = empty.And(rhs)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint64(0), result.Count())

	result, err = rhs.And(empty)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint64(0), result.Count())

	result, err = lhs.And(rhs)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, lhs.Count(), result.Count())
	assert.Equal(t, rhs.Count(), result.Count())
}

func Test_BitSet_AndCount(t *testing.T) {
	lhs := newBitSet(testCauldron)
	rhs := newBitSet(testCauldron)

	// And two empty sets
	count, err := lhs.AndCount(rhs)
	assert.Nil(t, err)
	assert.Equal(t, uint64(0), count)

	// One empty, the other with values
	empty := newBitSet(testCauldron)
	lhs = newBitSet(testCauldron)
	rhs = newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := uint64(rand.Int63n(TestSetSize))
		rhs.SetBit(bitIdx)
		lhs.SetBit(bitIdx)
	}

	count, err = empty.AndCount(rhs)
	assert.Nil(t, err)
	assert.Equal(t, uint64(0), count)

	count, err = rhs.AndCount(empty)
	assert.Nil(t, err)
	assert.Equal(t, uint64(0), count)

	count, err = lhs.AndCount(rhs)
	assert.Nil(t, err)
	assert.Equal(t, lhs.Count(), count)
	assert.Equal(t, rhs.Count(), count)
}

func Test_BitSet_Or(t *testing.T) {
	lhs := newBitSet(testCauldron)
	rhs := newBitSet(testCauldron)

	// Or two empty sets
	result, err := lhs.Or(rhs)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint64(0), result.Count())

	// One empty, the other with values
	empty := newBitSet(testCauldron)
	lhs = newBitSet(testCauldron)
	rhs = newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := uint64(rand.Int63n(TestSetSize))
		rhs.SetBit(bitIdx)
		lhs.SetBit(bitIdx)
	}

	result, err = empty.Or(rhs)
	assert.Nil(t, err, err)
	assert.NotNil(t, result)
	assert.Equal(t, rhs.Count(), result.Count())

	result, err = rhs.Or(empty)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, rhs.Count(), result.Count())

	result, err = lhs.Or(rhs)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, lhs.Count(), result.Count())
	assert.Equal(t, rhs.Count(), result.Count())
}

func Test_BitSet_OrCount(t *testing.T) {
	lhs := newBitSet(testCauldron)
	rhs := newBitSet(testCauldron)

	// Or two empty sets
	count, err := lhs.OrCount(rhs)
	assert.Nil(t, err)
	assert.Equal(t, uint64(0), count)

	// One empty, the other with values
	empty := newBitSet(testCauldron)
	lhs = newBitSet(testCauldron)
	rhs = newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := uint64(rand.Int63n(TestSetSize))
		rhs.SetBit(bitIdx)
		lhs.SetBit(bitIdx)
	}

	count, err = empty.OrCount(rhs)
	assert.Nil(t, err)
	assert.Equal(t, rhs.Count(), count)

	count, err = rhs.OrCount(empty)
	assert.Nil(t, err)
	assert.Equal(t, rhs.Count(), count)

	count, err = lhs.OrCount(rhs)
	assert.Nil(t, err)
	assert.Equal(t, lhs.Count(), count)
	assert.Equal(t, rhs.Count(), count)
}

func Test_BitSet_Not(t *testing.T) {
	set := newBitSet(testCauldron)
	cauldronCount := testCauldron.Count()

	// Not empty set
	result, err := set.Not()
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cauldronCount, result.Count())

	// Add some values
	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	set = newBitSet(testCauldron)
	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := uint64(rand.Int63n(TestSetSize))
		set.SetBit(bitIdx)
	}

	countBeforeNot := set.Count()

	result, err = set.Not()
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cauldronCount-countBeforeNot, result.Count())
}

func Test_BitSet_NotCount(t *testing.T) {
	set := newBitSet(testCauldron)
	cauldronCount := testCauldron.Count()

	// Not empty set
	count, err := set.NotCount()
	assert.Nil(t, err)
	assert.Equal(t, cauldronCount, count)

	// Add some values
	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	set = newBitSet(testCauldron)
	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := uint64(rand.Int63n(TestSetSize))
		set.SetBit(bitIdx)
	}

	countBeforeNot := set.Count()

	count, err = set.NotCount()
	assert.Nil(t, err)
	assert.Equal(t, cauldronCount-countBeforeNot, count)
}

func Test_BitSet_Size(t *testing.T) {
	set := newBitSet(testCauldron)

	assert.Equal(t, testCauldron.Capacity(), set.Size())
}

func Test_BitSet_ToSlice(t *testing.T) {
	set := newBitSet(testCauldron)
	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := uint64(rand.Int63n(TestSetSize))
		set.SetBit(bitIdx)
	}

	slice := set.ToSlice()
	assert.Equal(t, set.Count(), uint64(len(slice)))
}

func Test_BitSet_SetBit_GetBit_UnsetBit(t *testing.T) {
	set := newBitSet(testCauldron)
	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := uint64(rand.Int63n(TestSetSize))

		err := set.SetBit(bitIdx)
		assert.Nil(t, err)

		bit, err := set.GetBit(bitIdx)
		assert.Nil(t, err)
		assert.True(t, bit)

		err = set.UnsetBit(bitIdx)
		assert.Nil(t, err)

		bit, err = set.GetBit(bitIdx)
		assert.Nil(t, err)
		assert.False(t, bit)
	}
}

func Test_BitSet_Clear(t *testing.T) {
	set := newBitSet(testCauldron)
	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := uint64(rand.Int63n(TestSetSize))
		set.SetBit(bitIdx)
	}

	set.Clear()
	assert.Equal(t, uint64(0), set.Count(), set.Count())
}

func Test_BitSet_Fill(t *testing.T) {
	set := newBitSet(testCauldron)

	// fill empty set
	set.Fill()
	cap := testCauldron.Capacity()
	count := set.Count()
	assert.Equal(t, count, cap)

	set = newBitSet(testCauldron)
	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := uint64(rand.Int63n(TestSetSize))
		set.SetBit(bitIdx)
	}

	// fill partial set
	set.Fill()
	cap = testCauldron.Capacity()
	count = set.Count()
	assert.Equal(t, count, cap)
}

func Benchmark_BitSet_And(b *testing.B) {

	lhs := newBitSet(testCauldron)
	rhs := newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := uint64(rand.Int63n(TestSetSize))
		rhs.SetBit(bitIdx)
		bitIdx = uint64(rand.Int63n(TestSetSize))
		lhs.SetBit(bitIdx)
	}

	for n := 0; n < b.N; n++ {
		lhs.And(rhs)
	}
}

func Benchmark_BitSet_AndCount(b *testing.B) {

	lhs := newBitSet(testCauldron)
	rhs := newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := uint64(rand.Int63n(TestSetSize))
		rhs.SetBit(bitIdx)
		bitIdx = uint64(rand.Int63n(TestSetSize))
		lhs.SetBit(bitIdx)
	}

	for n := 0; n < b.N; n++ {
		lhs.AndCount(rhs)
	}
}

func Benchmark_BitSet_Or(b *testing.B) {

	lhs := newBitSet(testCauldron)
	rhs := newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := uint64(rand.Int63n(TestSetSize))
		rhs.SetBit(bitIdx)
		bitIdx = uint64(rand.Int63n(TestSetSize))
		lhs.SetBit(bitIdx)
	}

	for n := 0; n < b.N; n++ {
		lhs.Or(rhs)
	}
}

func Benchmark_BitSet_OrCount(b *testing.B) {

	lhs := newBitSet(testCauldron)
	rhs := newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := uint64(rand.Int63n(TestSetSize))
		rhs.SetBit(bitIdx)
		bitIdx = uint64(rand.Int63n(TestSetSize))
		lhs.SetBit(bitIdx)
	}

	for n := 0; n < b.N; n++ {
		lhs.OrCount(rhs)
	}
}


func Benchmark_BitSet_Not(b *testing.B) {

	set := newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := uint64(rand.Int63n(TestSetSize))
		set.SetBit(bitIdx)
	}

	for n := 0; n < b.N; n++ {
		set.Not()
	}
}

func Benchmark_BitSet_NotCount(b *testing.B) {

	set := newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := uint64(rand.Int63n(TestSetSize))
		set.SetBit(bitIdx)
	}

	for n := 0; n < b.N; n++ {
		set.NotCount()
	}
}
