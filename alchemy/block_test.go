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

	bb.fill(63)
	bit, err := bb.getBit(62)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = bb.getBit(63)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)

	bb.fill(65)
	bit, err = bb.getBit(64)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = bb.getBit(65)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)

	bb.fill(1021)
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

	bb.clear(35)
	bit, err = bb.getBit(34)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)
	bit, err = bb.getBit(35)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)
	bit, err = bb.getBit(36)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)

}

