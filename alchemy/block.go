package alchemy

import (
	"github.com/go-errors/errors"
	"math/bits"
)

const (
	BitsPerLong       = 64
	LongsPerBlock     = 16
	BitsPerBlock      = LongsPerBlock * BitsPerLong
	ZeroBitPattern    = Long(0)
	AllOnesBitPattern = Long(18446744073709551615)
)

type bitBlock struct {
	bits [LongsPerBlock]Long
}

func newBitBlock() *bitBlock {
	return &bitBlock{
	}
}

func newBitBlockWithBits(bits [LongsPerBlock]Long) *bitBlock {
	return &bitBlock{
		bits: bits,
	}
}

func (bb bitBlock) and(rhs *bitBlock) *bitBlock {
	var result [LongsPerBlock]Long

	for i := 0; i < LongsPerBlock; i++ {
		result[i] = bb.bits[i] & rhs.bits[i]
	}

	return newBitBlockWithBits(result)
}

func (bb bitBlock) or(rhs *bitBlock) *bitBlock {
	var result [LongsPerBlock]Long

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

	var result [LongsPerBlock]Long

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

func (bb bitBlock) andCount(rhs *bitBlock) Long {
	var count Long = 0

	for i := 0; i < LongsPerBlock; i++ {
		result := bb.bits[i] & rhs.bits[i]
		count += bb.numberOfSetBits(result)
	}

	return count
}

func (bb bitBlock) orCount(rhs *bitBlock) Long {
	var count Long = 0

	for i := 0; i < LongsPerBlock; i++ {
		result := bb.bits[i] | rhs.bits[i]
		count += bb.numberOfSetBits(result)
	}

	return count
}

func (bb bitBlock) notAllCount() Long {
	count, _ := bb.notCount(BitsPerBlock)
	return count
}

func (bb bitBlock) notCount(len int) (Long, error) {
	if len > BitsPerBlock {
		return ZeroBitPattern, errors.New("length out of range for NotCount(len int)")
	}

	var count Long = 0

	numLongs := len / BitsPerLong
	lastBits := len % BitsPerLong

	for i := 0; i < numLongs; i++ {
		result := ^bb.bits[i]
		count += bb.numberOfSetBits(result)
	}

	// create a mask And apply the last one
	if lastBits > 0 {
		lastLong := numLongs
		mask := bb.getMask(lastBits)
		notBits := ^(bb.bits[lastLong])
		notBits &= mask
		count += bb.numberOfSetBits(notBits)
	}

	return count, nil
}

func (bb bitBlock) countAll() Long {
	count, _ := bb.count(BitsPerBlock)
	return count
}

func (bb bitBlock) count(len int) (Long, error) {
	if len > BitsPerBlock {
		return ZeroBitPattern, errors.New("length out of range for count(len int)")
	}

	var count Long = 0

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
	shift := BitsPerLong - Long(1) - Long(bit)
	mask := Long(1) << shift

	return (bits & mask) != 0, nil
}

func (bb *bitBlock) setBit(index int) error {
	if index >= BitsPerBlock {
		return errors.New("length out of range for SetBit(index int)")
	}

	i := index / BitsPerLong
	bit := index % BitsPerLong
	shift := BitsPerLong - Long(1) - Long(bit)
	mask := Long(1) << shift

	bb.bits[i] |= mask
	return nil
}

func (bb *bitBlock) unsetBit(index int) error {
	if index >= BitsPerBlock {
		return errors.New("length out of range for UnsetBit(index int)")
	}

	i := index / BitsPerLong
	bit := index % BitsPerLong
	mask := Long(1) << (BitsPerLong - Long(1) - Long(bit))

	bb.bits[i] &= ^mask
	return nil
}

func (bb *bitBlock) fillAll() {
	bb.fill(BitsPerBlock)
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

	if lastBits > 0 {
		lastLong := numLongs
		mask := bb.getMask(lastBits)
		bb.bits[lastLong] = mask
	}

	return nil
}

func (bb *bitBlock) clearAll() {
	bb.clear(BitsPerBlock)
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

func (bb bitBlock) getMask(bitIndex int) Long {
	var shift Long
	var l Long

	shift = Long(BitsPerLong - bitIndex)
	l = AllOnesBitPattern
	return l << shift
}

func (bb bitBlock) numberOfSetBits(x Long) Long {
	return Long(bits.OnesCount64(uint64(x)))
}

