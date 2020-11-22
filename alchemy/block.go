package alchemy

const (
	BitsPerLong   = 64
	LongsPerBlock = 16
	BitsPerBlock  = LongsPerBlock * BitsPerLong
	ZeroBitPattern = 0
	LargestLong = 18446744073709551615
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
	return bb.not(BitsPerBlock)
}

func (bb *bitBLock) not(len int) *bitBLock {

	var result [LongsPerBlock]Long

	numLongs := len / BitsPerLong
	lastBits := len % BitsPerLong

	for i := 0; i < numLongs; i++ {
		result[i] = ^bb.bits[i]
	}

	// create a mask and apply the last one
	if lastBits > 0 {
		lastLong := numLongs
		if (LongsPerBlock - 1) < numLongs {
			lastLong = LongsPerBlock - 1
		}

		mask := bb.getMask(lastBits)
		notBits := ^(bb.bits[lastLong])
		notBits &= mask
		result[lastLong] = notBits
	}

	return newBitBlockWithBits(result)
}

func (bb *bitBLock)  getMask(bitIndex int) Long {
	var shift Long
	var l Long

	shift = Long(BitsPerLong - bitIndex)
	l = LargestLong
	return l << shift
}