package alchemy

import (
	"testing"

	"gotest.tools/assert"
)

func Test_New_BitBlock(t *testing.T) {
	bb := newBitBlock()
	assert.Assert(t, bb != nil, bb)
	assert.Assert(t, bb.bits[0] == 0, bb.bits[0])

	var bits [LongsPerBlock]Long
	bits[0] = AllOnesBitPattern

	bb = newBitBlockWithBits(bits)
	assert.Assert(t, bb != nil, bb)
	assert.Assert(t, bb.bits[0] == AllOnesBitPattern, bb.bits[0])
}

func Test_GetBit(t *testing.T) {
	var bits [LongsPerBlock]Long
	bits[0] = AllOnesBitPattern

	bb := newBitBlockWithBits(bits)

	// first 64 bits should all be 1s
	for i := 0; i < BitsPerLong; i++ {
		bit, err := bb.getBit(i)
		assert.Assert(t, err == nil, err)
		assert.Assert(t, bit, bit)
	}

	// next 64 bits should all be 0s
	for i := BitsPerLong; i < (BitsPerLong+BitsPerLong); i++ {
		bit, err := bb.getBit(i)
		assert.Assert(t, err == nil, err)
		assert.Assert(t, !bit, bit)
	}

	// out of range
	bit, err := bb.getBit(BitsPerBlock)
	assert.Assert(t, err != nil, err)
	assert.Assert(t, !bit, bit)
}

func Test_Set_Unset_Bit(t *testing.T) {
	bb := newBitBlock()

	bit, err := bb.getBit(0)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)

	err = bb.unsetBit(0)
	assert.Assert(t, err == nil, err)
	bit, err = bb.getBit(0)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)

	err = bb.setBit(0)
	assert.Assert(t, err == nil, err)
	bit, err = bb.getBit(0)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)

	err = bb.setBit(1023)
	assert.Assert(t, err == nil, err)
	bit, err = bb.getBit(1023)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)

	err = bb.unsetBit(0)
	assert.Assert(t, err == nil, err)
	bit, err = bb.getBit(0)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)

	// out of range
	err = bb.setBit(1024)
	assert.Assert(t, err != nil, err)
	err = bb.unsetBit(1024)
	assert.Assert(t, err != nil, err)
}

func Test_Fill_Clear(t *testing.T) {
	bb := newBitBlock()

	bb.fillAll()
	for i := 0; i < BitsPerBlock; i++ {
		bit, err := bb.getBit(i)
		assert.Assert(t, err == nil, err)
		assert.Assert(t, bit, bit)
	}

	bb.clearAll()
	for i := 0; i < BitsPerBlock; i++ {
		bit, err := bb.getBit(i)
		assert.Assert(t, err == nil, err)
		assert.Assert(t, !bit, bit)
	}

	err := bb.fill(63)
	assert.Assert(t, err == nil, err)
	bit, err := bb.getBit(62)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = bb.getBit(63)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)

	err = bb.fill(65)
	assert.Assert(t, err == nil, err)
	bit, err = bb.getBit(64)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = bb.getBit(65)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)

	err = bb.fill(1021)
	assert.Assert(t, err == nil, err)
	bit, err = bb.getBit(1020)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = bb.getBit(1021)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)
	bit, err = bb.getBit(1022)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)
	bit, err = bb.getBit(1023)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)

	err = bb.clear(35)
	assert.Assert(t, err == nil, err)
	bit, err = bb.getBit(34)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)
	bit, err = bb.getBit(35)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = bb.getBit(36)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)

	// out of range
	err = bb.fill(1025)
	assert.Assert(t, err != nil, err)
	err = bb.clear(1025)
	assert.Assert(t, err != nil, err)
}

func Test_And_AndCount(t *testing.T) {
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
	assert.Assert(t, count==6, count)
	count = lhs.andCount(rhs)
	assert.Assert(t, count==6, count)

	bit, err := result.getBit(0)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = result.getBit(1)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = result.getBit(2)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = result.getBit(1021)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = result.getBit(1022)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = result.getBit(1023)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
}

func Test_Or_OrCount(t *testing.T) {
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
	assert.Assert(t, count==10, count)
	count = lhs.orCount(rhs)
	assert.Assert(t, count==10, count)

	bit, err := result.getBit(0)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = result.getBit(1)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = result.getBit(2)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = result.getBit(3)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = result.getBit(4)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = result.getBit(1019)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = result.getBit(1020)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = result.getBit(1021)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = result.getBit(1022)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = result.getBit(1023)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
}

func Test_Not_NotCount(t *testing.T) {
	bb := newBitBlock()

	bb.setBit(0)
	bb.setBit(1)
	bb.setBit(2)
	bb.setBit(1021)
	bb.setBit(1022)
	bb.setBit(1023)

	count := bb.notAllCount()
	assert.Assert(t, count == 1018, count)

	result := bb.notAll()
	bit, err := result.getBit(0)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)
	bit, err = result.getBit(1)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)
	bit, err = result.getBit(2)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)
	bit, err = result.getBit(1021)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)
	bit, err = result.getBit(1022)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)
	bit, err = result.getBit(1023)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)

	result, err = bb.not(510)
	assert.Assert(t, err == nil, err)
	count = result.countAll()
	assert.Assert(t, count==507, count)

	count, err = bb.notCount(510)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, count==507, count)

	result, err = bb.not(1025)
	assert.Assert(t, err != nil, err)

	count, err = bb.notCount(1026)
	assert.Assert(t, err != nil, err)
}

func Test_Count(t *testing.T) {
	bb := newBitBlock()

	count := bb.countAll()
	assert.Assert(t, count==0, count)

	bb.setBit(0)
	bb.setBit(1)
	bb.setBit(2)
	bb.setBit(1021)
	bb.setBit(1022)
	bb.setBit(1023)

	count, err := bb.count(10)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, count==3, count)

	count, err = bb.count(1025)
	assert.Assert(t, err != nil, err)
}

func Benchmark_And(b *testing.B) {

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

func Benchmark_AndCount(b *testing.B) {

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