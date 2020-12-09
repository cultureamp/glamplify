package alchemy

import (
	"math/bits"

	"github.com/go-errors/errors"
)

const (
	// BitsPerLong = 64
	BitsPerLong       = 64
	// LongsPerBlock = 16
	LongsPerBlock     = 16
	// BitsPerBlock = LongsPerBlock * BitsPerLong
	BitsPerBlock      = LongsPerBlock * BitsPerLong
	// ZeroBitPattern = uint64(0)
	ZeroBitPattern    = uint64(0)
	// AllOnesBitPattern = uint64(18446744073709551615)
	AllOnesBitPattern = uint64(18446744073709551615)
)

type bitBlock struct {
	bits [LongsPerBlock]uint64
}

func newBitBlock() *bitBlock {
	return &bitBlock{
	}
}

func newBitBlockWithBits(bits [LongsPerBlock]uint64) *bitBlock {
	return &bitBlock{
		bits: bits,
	}
}

func (bb bitBlock) and(rhs *bitBlock) *bitBlock {
	var result [LongsPerBlock]uint64

	for i := 0; i < LongsPerBlock; i++ {
		result[i] = bb.bits[i] & rhs.bits[i]
	}

	return newBitBlockWithBits(result)
}

func (bb bitBlock) or(rhs *bitBlock) *bitBlock {
	var result [LongsPerBlock]uint64

	for i := 0; i < LongsPerBlock; i++ {
		result[i] = bb.bits[i] | rhs.bits[i]
	}

	return newBitBlockWithBits(result)
}

func (bb *bitBlock) notAll() *bitBlock {
	block, _ := bb.not(BitsPerBlock)
	return block
}

func (bb bitBlock) not(len int) (*bitBlock, error) {
	if len > BitsPerBlock {
		return nil, errors.New("length out of range for Not(len int)")
	}

	var result [LongsPerBlock]uint64

	numLongs := len / BitsPerLong
	lastBits := len % BitsPerLong

	for i := 0; i < numLongs; i++ {
		result[i] = ^bb.bits[i]
	}

	// create a mask And apply the last one
	if lastBits > 0 {
		lastLong := numLongs
		mask := bb.getMask(lastBits)
		notBits := ^(bb.bits[lastLong])
		notBits &= mask
		result[lastLong] = notBits
	}

	return newBitBlockWithBits(result), nil
}

func (bb bitBlock) andCount(rhs *bitBlock) uint64 {
	var count uint64 = 0

	for i := 0; i < LongsPerBlock; i++ {
		result := bb.bits[i] & rhs.bits[i]
		count += bb.numberOfSetBits(result)
	}

	return count
}

func (bb bitBlock) orCount(rhs *bitBlock) uint64 {
	var count uint64 = 0

	for i := 0; i < LongsPerBlock; i++ {
		result := bb.bits[i] | rhs.bits[i]
		count += bb.numberOfSetBits(result)
	}

	return count
}

func (bb bitBlock) countAll() uint64 {
	count, _ := bb.count(BitsPerBlock)
	return count
}

func (bb bitBlock) count(len int) (uint64, error) {
	if len > BitsPerBlock {
		return ZeroBitPattern, errors.New("length out of range for count(len int)")
	}

	var count uint64 = 0

	numLongs := len / BitsPerLong
	lastBits := len % BitsPerLong

	for i := 0; i < numLongs; i++ {
		count += bb.numberOfSetBits(bb.bits[i])
	}

	// create a mask And apply the last one
	if lastBits > 0 {
		lastLong := numLongs
		mask := bb.getMask(lastBits)
		bits := bb.bits[lastLong]
		bits &= mask
		count += bb.numberOfSetBits(bits)
	}

	return count, nil
}

func (bb bitBlock) getBit(index int) (bool, error) {
	if index >= BitsPerBlock {
		return false, errors.New("length out of range for GetBit(index int)")
	}
	// http://stackoverflow.com/questions/4854207/get-a-specific-bit-from-byte

	i := index / BitsPerLong
	bits := bb.bits[i]
	bit := index % BitsPerLong
	shift := BitsPerLong - 1 - bit
	mask := uint64(1) << shift

	return (bits & mask) != 0, nil
}

func (bb *bitBlock) setBit(index int) error {
	if index >= BitsPerBlock {
		return errors.New("length out of range for SetBit(index int)")
	}

	i := index / BitsPerLong
	bit := index % BitsPerLong
	shift := BitsPerLong - 1 - bit
	mask := uint64(1) << shift

	bb.bits[i] |= mask
	return nil
}

func (bb *bitBlock) unsetBit(index int) error {
	if index >= BitsPerBlock {
		return errors.New("length out of range for UnsetBit(index int)")
	}

	i := index / BitsPerLong
	bit := index % BitsPerLong
	mask := uint64(1) << (BitsPerLong - 1 - bit)

	bb.bits[i] &= ^mask
	return nil
}

func (bb *bitBlock) fillAll() error {
	return bb.fill(BitsPerBlock)
}

func (bb *bitBlock) fill(len int) error {
	if len > BitsPerBlock {
		return errors.New("length out of range for Fill(len int)")
	}

	numLongs := len / BitsPerLong
	lastBits := len % BitsPerLong

	for i := 0; i < numLongs; i++ {
		bb.bits[i] = AllOnesBitPattern
	}

	lastLong := numLongs
	if lastBits > 0 {
		mask := bb.getMask(lastBits)
		bb.bits[lastLong] = mask
		lastLong++
	}

	for i := lastLong; i < LongsPerBlock; i++ {
		bb.bits[i] = ZeroBitPattern
	}

	return nil
}

func (bb *bitBlock) clearAll() error {
	return bb.clear(BitsPerBlock)
}

func (bb *bitBlock) clear(len int) error {
	if len > BitsPerBlock {
		return errors.New("length out of range for Fill(len int)")
	}

	numLongs := len / BitsPerLong
	lastBits := len % BitsPerLong

	for i := 0; i < numLongs; i++ {
		bb.bits[i] = ZeroBitPattern
	}

	if lastBits > 0 {
		lastLong := numLongs
		mask := bb.getMask(lastBits)
		bb.bits[lastLong] &= ^mask
	}

	return nil
}

func (bb bitBlock) getMask(bitIndex int) uint64 {
	var l uint64

	shift := BitsPerLong - bitIndex
	l = AllOnesBitPattern
	return l << shift
}

func (bb bitBlock) numberOfSetBits(x uint64) uint64 {
	return uint64(bits.OnesCount64(x))
}

