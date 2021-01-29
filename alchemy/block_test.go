package alchemy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New_BitBlock(t *testing.T) {
	bb := newBitBlock()
	assert.NotNil(t, bb)
	assert.Equal(t, uint64(0), bb.bits[0])

	var bits [LongsPerBlock]uint64
	bits[0] = AllOnesBitPattern

	bb = newBitBlockWithBits(bits)
	assert.NotNil(t, bb)
	assert.Equal(t, AllOnesBitPattern, bb.bits[0])
}

func Test_BitBlock_GetBit(t *testing.T) {
	var bits [LongsPerBlock]uint64
	bits[0] = AllOnesBitPattern

	bb := newBitBlockWithBits(bits)

	// first 64 bits should all be 1s
	for i := 0; i < BitsPerLong; i++ {
		bit, err := bb.getBit(i)
		assert.Nil(t, err)
		assert.True(t, bit)
	}

	// next 64 bits should all be 0s
	for i := BitsPerLong; i < (BitsPerLong+BitsPerLong); i++ {
		bit, err := bb.getBit(i)
		assert.Nil(t, err)
		assert.False(t, bit)
	}

	// out of range
	bit, err := bb.getBit(BitsPerBlock)
	assert.NotNil(t, err)
	assert.False(t, bit)
}

func Test_BitBlock_Set_Unset_Bit(t *testing.T) {
	bb := newBitBlock()

	bit, err := bb.getBit(0)
	assert.Nil(t, err)
	assert.False(t, bit)

	err = bb.unsetBit(0)
	assert.Nil(t, err)
	bit, err = bb.getBit(0)
	assert.Nil(t, err)
	assert.False(t, bit)

	err = bb.setBit(0)
	assert.Nil(t, err)
	bit, err = bb.getBit(0)
	assert.Nil(t, err)
	assert.True(t, bit)

	err = bb.setBit(1023)
	assert.Nil(t, err)
	bit, err = bb.getBit(1023)
	assert.Nil(t, err)
	assert.True(t, bit)

	err = bb.unsetBit(0)
	assert.Nil(t, err)
	bit, err = bb.getBit(0)
	assert.Nil(t, err)
	assert.False(t, bit)

	// out of range
	err = bb.setBit(1024)
	assert.NotNil(t, err)
	err = bb.unsetBit(1024)
	assert.NotNil(t, err)
}

func Test_BitBlock_Fill_Clear(t *testing.T) {
	bb := newBitBlock()

	bb.fillAll()
	for i := 0; i < BitsPerBlock; i++ {
		bit, err := bb.getBit(i)
		assert.Nil(t, err)
		assert.True(t, bit)
	}

	bb.clearAll()
	for i := 0; i < BitsPerBlock; i++ {
		bit, err := bb.getBit(i)
		assert.Nil(t, err)
		assert.False(t, bit)
	}

	bb.fill(63)
	bit, err := bb.getBit(62)
	assert.Nil(t, err)
	assert.True(t, bit)
	bit, err = bb.getBit(63)
	assert.Nil(t, err)
	assert.False(t, bit)

	bb.fill(65)
	assert.Nil(t, err)
	bit, err = bb.getBit(64)
	assert.Nil(t, err)
	assert.True(t, bit)
	bit, err = bb.getBit(65)
	assert.Nil(t, err)
	assert.False(t, bit)

	bb.fill(1021)
	assert.Nil(t, err)
	bit, err = bb.getBit(1020)
	assert.Nil(t, err)
	assert.True(t, bit)
	bit, err = bb.getBit(1021)
	assert.Nil(t, err)
	assert.False(t, bit)
	bit, err = bb.getBit(1022)
	assert.Nil(t, err)
	assert.False(t, bit)
	bit, err = bb.getBit(1023)
	assert.Nil(t, err)
	assert.False(t, bit)

	bb.clear(35)
	assert.Nil(t, err)
	bit, err = bb.getBit(34)
	assert.Nil(t, err)
	assert.False(t, bit)
	bit, err = bb.getBit(35)
	assert.Nil(t, err)
	assert.True(t, bit)
	bit, err = bb.getBit(36)
	assert.Nil(t, err)
	assert.True(t, bit)

	// out of range
	bb.fill(1025)
	bit, err = bb.getBit(1023)
	assert.Nil(t, err)
	assert.True(t, bit)
	bb.clear(1025)
	bit, err = bb.getBit(1023)
	assert.Nil(t, err)
	assert.False(t, bit)
}

func Test_BitBlock_And_AndCount(t *testing.T) {
	lhs := newBitBlock()
	rhs := newBitBlock()

	lhs.setBit(0)
	lhs.setBit(1)
	lhs.setBit(2)
	lhs.setBit(1021)
	lhs.setBit(1022)
	lhs.setBit(1023)

	rhs.setBit(0)
	rhs.setBit(1)
	rhs.setBit(2)
	rhs.setBit(1021)
	rhs.setBit(1022)
	rhs.setBit(1023)

	result := lhs.and(rhs)
	count := result.countAll()
	assert.Equal(t, uint64(6), count)
	count = lhs.andCount(rhs)
	assert.Equal(t, uint64(6), count)

	bit, err := result.getBit(0)
	assert.Nil(t, err)
	assert.True(t, bit)
	bit, err = result.getBit(1)
	assert.Nil(t, err)
	assert.True(t, bit)
	bit, err = result.getBit(2)
	assert.Nil(t, err)
	assert.True(t, bit)
	bit, err = result.getBit(1021)
	assert.Nil(t, err)
	assert.True(t, bit)
	bit, err = result.getBit(1022)
	assert.Nil(t, err)
	assert.True(t, bit)
	bit, err = result.getBit(1023)
	assert.Nil(t, err)
	assert.True(t, bit)
}

func Test_BitBlock_Or_OrCount(t *testing.T) {
	lhs := newBitBlock()
	rhs := newBitBlock()

	lhs.setBit(0)
	lhs.setBit(1)
	lhs.setBit(2)
	lhs.setBit(1021)
	lhs.setBit(1022)
	lhs.setBit(1023)

	rhs.setBit(2)
	rhs.setBit(3)
	rhs.setBit(4)
	rhs.setBit(1019)
	rhs.setBit(1020)
	rhs.setBit(1021)

	result := lhs.or(rhs)
	count := result.countAll()
	assert.Equal(t, uint64(10), count)
	count = lhs.orCount(rhs)
	assert.Equal(t, uint64(10), count)

	bit, err := result.getBit(0)
	assert.Nil(t, err)
	assert.True(t, bit)
	bit, err = result.getBit(1)
	assert.Nil(t, err)
	assert.True(t, bit)
	bit, err = result.getBit(2)
	assert.Nil(t, err)
	assert.True(t, bit)
	bit, err = result.getBit(3)
	assert.Nil(t, err)
	assert.True(t, bit)
	bit, err = result.getBit(4)
	assert.Nil(t, err)
	assert.True(t, bit)
	bit, err = result.getBit(1019)
	assert.Nil(t, err)
	assert.True(t, bit)
	bit, err = result.getBit(1020)
	assert.Nil(t, err)
	assert.True(t, bit)
	bit, err = result.getBit(1021)
	assert.Nil(t, err)
	assert.True(t, bit)
	bit, err = result.getBit(1022)
	assert.Nil(t, err)
	assert.True(t, bit)
	bit, err = result.getBit(1023)
	assert.Nil(t, err)
	assert.True(t, bit)
}

func Test_BitBlock_Count(t *testing.T) {
	bb := newBitBlock()

	count := bb.countAll()
	assert.Equal(t, uint64(0), count)

	bb.setBit(0)
	bb.setBit(1)
	bb.setBit(2)
	bb.setBit(1021)
	bb.setBit(1022)
	bb.setBit(1023)

	count, err := bb.count(10)
	assert.Nil(t, err)
	assert.Equal(t, uint64(3), count)

	count, err = bb.count(1025)
	assert.NotNil(t, err)
}

func Benchmark_BitBlock_And(b *testing.B) {
	lhs := newBitBlock()
	rhs := newBitBlock()

	lhs.setBit(0)
	lhs.setBit(1)
	lhs.setBit(2)
	lhs.setBit(1021)
	lhs.setBit(1022)
	lhs.setBit(1023)

	rhs.setBit(0)
	rhs.setBit(1)
	rhs.setBit(2)
	rhs.setBit(1021)
	rhs.setBit(1022)
	rhs.setBit(1023)

	for n := 0; n < b.N; n++ {
		lhs.and(rhs)
	}
}

func Benchmark_BitBlock_AndCount(b *testing.B) {
	lhs := newBitBlock()
	rhs := newBitBlock()

	lhs.setBit(0)
	lhs.setBit(1)
	lhs.setBit(2)
	lhs.setBit(1021)
	lhs.setBit(1022)
	lhs.setBit(1023)

	rhs.setBit(0)
	rhs.setBit(1)
	rhs.setBit(2)
	rhs.setBit(1021)
	rhs.setBit(1022)
	rhs.setBit(1023)

	for n := 0; n < b.N; n++ {
		lhs.andCount(rhs)
	}
}

func Benchmark_BitBlock_Or(b *testing.B) {
	lhs := newBitBlock()
	rhs := newBitBlock()

	lhs.setBit(0)
	lhs.setBit(1)
	lhs.setBit(2)
	lhs.setBit(1021)
	lhs.setBit(1022)
	lhs.setBit(1023)

	rhs.setBit(0)
	rhs.setBit(1)
	rhs.setBit(2)
	rhs.setBit(3)
	rhs.setBit(511)
	rhs.setBit(512)

	for n := 0; n < b.N; n++ {
		lhs.or(rhs)
	}
}

func Benchmark_BitBlock_OrCount(b *testing.B) {
	lhs := newBitBlock()
	rhs := newBitBlock()

	lhs.setBit(0)
	lhs.setBit(1)
	lhs.setBit(2)
	lhs.setBit(1021)
	lhs.setBit(1022)
	lhs.setBit(1023)

	rhs.setBit(0)
	rhs.setBit(1)
	rhs.setBit(2)
	rhs.setBit(3)
	rhs.setBit(511)
	rhs.setBit(512)

	for n := 0; n < b.N; n++ {
		lhs.orCount(rhs)
	}
}

func Benchmark_BitBlock_Not(b *testing.B) {
	set := newBitBlock()

	set.setBit(0)
	set.setBit(1)
	set.setBit(2)
	set.setBit(1021)
	set.setBit(1022)
	set.setBit(1023)

	for n := 0; n < b.N; n++ {
		set.notAll()
	}
}
