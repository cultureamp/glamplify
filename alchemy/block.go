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

type bitBLock struct {
	bits [LongsPerBlock]Long
}

func newBitBlock() *bitBLock {
	return &bitBLock{
	}
}

func newBitBlockWithBits(bits [LongsPerBlock]Long) *bitBLock {
	return &bitBLock{
		bits: bits,
	}
}

func (bb *bitBLock) and(rhs *bitBLock) *bitBLock {
	var result [LongsPerBlock]Long

	for i := 0; i < LongsPerBlock; i++ {
		result[i] = bb.bits[i] & rhs.bits[i]
	}

	return newBitBlockWithBits(result)
}

func (bb *bitBLock) or(rhs *bitBLock) *bitBLock {
	var result [LongsPerBlock]Long

	for i := 0; i < LongsPerBlock; i++ {
		result[i] = bb.bits[i] | rhs.bits[i]
	}

	return newBitBlockWithBits(result)
}

func (bb *bitBLock) notAll() *bitBLock {
	block, _ := bb.not(BitsPerBlock)
	return block
}

func (bb *bitBLock) not(len int) (*bitBLock, error) {
	if len > BitsPerBlock {
		return nil, errors.New("length out of range for not(len int)")
	}

	var result [LongsPerBlock]Long

	numLongs := len / BitsPerLong
	lastBits := len % BitsPerLong

	for i := 0; i < numLongs; i++ {
		result[i] = ^bb.bits[i]
	}

	// create a mask and apply the last one
	if lastBits > 0 {
		lastLong := numLongs
		//if (LongsPerBlock - 1) < numLongs {
		//	lastLong = LongsPerBlock - 1
		//}

		mask := bb.getMask(lastBits)
		notBits := ^(bb.bits[lastLong])
		notBits &= mask
		result[lastLong] = notBits
	}

	return newBitBlockWithBits(result), nil
}

func (bb *bitBLock) andCount(rhs *bitBLock) Long {
	var count Long = 0

	for i := 0; i < LongsPerBlock; i++ {
		result := bb.bits[i] & rhs.bits[i]
		count += bb.numberOfSetBits(result)
	}

	return count
}

func (bb *bitBLock) orCount(rhs *bitBLock) Long {
	var count Long = 0

	for i := 0; i < LongsPerBlock; i++ {
		result := bb.bits[i] | rhs.bits[i]
		count += bb.numberOfSetBits(result)
	}

	return count
}

func (bb *bitBLock) notAllCount() Long {
	count, _ := bb.notCount(BitsPerBlock)
	return count
}

func (bb *bitBLock) notCount(len int) (Long, error) {
	if len > BitsPerBlock {
		return ZeroBitPattern, errors.New("length out of range for notCount(len int)")
	}

	var count Long = 0

	numLongs := len / BitsPerLong
	lastBits := len % BitsPerLong

	for i := 0; i < numLongs; i++ {
		result := ^bb.bits[i]
		count += bb.numberOfSetBits(result)
	}

	// create a mask and apply the last one
	if lastBits > 0 {
		lastLong := numLongs
		//if (LongsPerBlock - 1) < numLongs {
		//	lastLong = LongsPerBlock - 1
		//}

		mask := bb.getMask(lastBits)
		notBits := ^(bb.bits[lastLong])
		notBits &= mask
		count += bb.numberOfSetBits(notBits)
	}

	return count, nil
}

func (bb *bitBLock) countAll() Long {
	count, _ := bb.count(BitsPerBlock)
	return count
}

func (bb *bitBLock) count(len int) (Long, error) {
	if len > BitsPerBlock {
		return ZeroBitPattern, errors.New("length out of range for count(len int)")
	}

	var count Long = 0

	numLongs := len / BitsPerLong
	lastBits := len % BitsPerLong

	for i := 0; i < numLongs; i++ {
		count += bb.numberOfSetBits(bb.bits[i])
	}

	// create a mask and apply the last one
	if lastBits > 0 {
		lastLong := numLongs
		//if (LongsPerBlock - 1) < numLongs {
		//	lastLong = LongsPerBlock - 1
		//}

		mask := bb.getMask(lastBits)
		bits := bb.bits[lastLong]
		bits &= mask
		count += bb.numberOfSetBits(bits)
	}

	return count, nil
}

func (bb *bitBLock) getBit(index int) (bool, error) {
	if index >= BitsPerBlock {
		return false, errors.New("length out of range for getBit(index int)")
	}
	// http://stackoverflow.com/questions/4854207/get-a-specific-bit-from-byte

	i := index / BitsPerLong
	bits := bb.bits[i]
	bit := index % BitsPerLong
	mask := Long(1) << (BitsPerLong - Long(1) - Long(bit))

	return (bits & mask) != 0, nil
}

func (bb *bitBLock) setBit(index int) error {
	if index >= BitsPerBlock {
		return errors.New("length out of range for setBit(index int)")
	}

	i := index / BitsPerLong
	bit := index % BitsPerLong
	mask := Long(1) << (BitsPerLong - Long(1) - Long(bit))

	bb.bits[i] |= mask
	return nil
}

func (bb *bitBLock) unsetBit(index int) error {
	if index >= BitsPerBlock {
		return errors.New("length out of range for unsetBit(index int)")
	}

	i := index / BitsPerLong
	bit := index % BitsPerLong
	mask := Long(1) << (BitsPerLong - Long(1) - Long(bit))

	bb.bits[i] &= ^mask
	return nil
}

func (bb *bitBLock) fillAll() {
	bb.fill(BitsPerBlock)
}

func (bb *bitBLock) fill(len int) error {
	if len > BitsPerBlock {
		return errors.New("length out of range for fill(len int)")
	}

	numLongs := len / BitsPerLong
	lastBits := len % BitsPerLong

	for i := 0; i < numLongs; i++ {
		bb.bits[i] = AllOnesBitPattern
	}

	if lastBits > 0 {
		lastLong := numLongs
		//if (LongsPerBlock - 1) < numLongs {
		//	lastLong = LongsPerBlock - 1
		//}

		mask := bb.getMask(lastBits)
		bb.bits[lastLong] = mask
	}

	return nil
}

func (bb *bitBLock) clearAll() {
	bb.clear(BitsPerBlock)
}

func (bb *bitBLock) clear(len int) error {
	if len > BitsPerBlock {
		return errors.New("length out of range for fill(len int)")
	}

	numLongs := len / BitsPerLong
	lastBits := len % BitsPerLong

	for i := 0; i < numLongs; i++ {
		bb.bits[i] = ZeroBitPattern
	}

	if lastBits > 0 {
		lastLong := numLongs
		//if (LongsPerBlock - 1) < numLongs {
		//	lastLong = LongsPerBlock - 1
		//}

		mask := bb.getMask(lastBits)
		bb.bits[lastLong] &= ^mask
	}

	return nil
}

func (bb *bitBLock) clone() *bitBLock {
	// copy (not * to)
	bits := bb.bits

	return newBitBlockWithBits(bits)
}

func (bb *bitBLock) getMask(bitIndex int) Long {
	var shift Long
	var l Long

	shift = Long(BitsPerLong - bitIndex)
	l = AllOnesBitPattern
	return l << shift
}

func (bb *bitBLock) numberOfSetBits(x Long) Long {
	return Long(bits.OnesCount64(uint64(x)))
}
