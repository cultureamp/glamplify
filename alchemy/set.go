package alchemy

import (
	"sync"

	"github.com/go-errors/errors"
)

// ReadOnlySet interface represents a set of bits that can not be modified
type ReadOnlySet interface {
	And(set ReadOnlySet) (Set, error)
	Or(set ReadOnlySet) (Set, error)
	Not() (Set, error)

	AndCount(set ReadOnlySet) (uint64, error)
	OrCount(set ReadOnlySet) (uint64, error)
	NotCount() (uint64, error)

	Count() uint64
	Size() uint64
	ToSlice() []Item

	GetBit(index uint64) (bool, error)
}

// Set interface represents a set of bits that can be modified
type Set interface {
	And(set ReadOnlySet) (Set, error)
	Or(set ReadOnlySet) (Set, error)
	Not() (Set, error)

	AndCount(set ReadOnlySet) (uint64, error)
	OrCount(set ReadOnlySet) (uint64, error)
	NotCount() (uint64, error)

	Count() uint64
	Size() uint64
	ToSlice() []Item

	GetBit(index uint64) (bool, error)
	SetBit(index uint64) error
	UnsetBit(index uint64) error
	Clear()
	Fill()
}

type bitSet struct {
	cauldron Cauldron
	blocks   map[int]*bitBlock
	lock     *sync.RWMutex
}

func newBitSet(cauldron Cauldron) Set {
	return &bitSet{
		cauldron: cauldron,
		blocks:   map[int]*bitBlock{},
		lock:     &sync.RWMutex{},
	}
}

func (set *bitSet) And(rhsSet ReadOnlySet) (Set, error) {
	rhs, ok := rhsSet.(*bitSet)
	if !ok {
		return nil, errors.New("invalid set passed to And")
	}

	set.lock.RLock()
	defer set.lock.RUnlock()

	result := newBitSet(set.cauldron).(*bitSet)

	blocks := set.getBlockCount()
	for i := 0; i < blocks; i++ {
		lhsBlock, lhsOk := set.getBlock(i)
		rhsBlock, rhsOk := rhs.getBlock(i)

		if !lhsOk || !rhsOk {
			// dont need to And if either left or right are empty
			continue
		}

		block := lhsBlock.and(rhsBlock)
		result.blocks[i] = block
	}

	return result, nil
}

func (set bitSet) Or(rhsSet ReadOnlySet) (Set, error) {
	rhs, ok := rhsSet.(*bitSet)
	if !ok {
		return nil, errors.New("invalid set passed to Or")
	}

	set.lock.RLock()
	defer set.lock.RUnlock()

	result := newBitSet(set.cauldron).(*bitSet)

	blocks := set.getBlockCount()
	for i := 0; i < blocks; i++ {
		lhsBlock, lhsOk := set.getBlock(i)
		rhsBlock, rhsOk := rhs.getBlock(i)

		var block *bitBlock
		if !lhsOk && !rhsOk {
			// dont need to Or if both left or right are empty
			continue
		} else if !rhsOk {
			block = lhsBlock
		} else if !lhsOk {
			block = rhsBlock
		} else {
			block = lhsBlock.or(rhsBlock)
		}

		result.blocks[i] = block
	}

	return result, nil
}

func (set bitSet) Not() (Set, error) {
	set.lock.RLock()
	defer set.lock.RUnlock()

	result := newBitSet(set.cauldron).(*bitSet)

	blocks := set.getBlockCount()
	for i := 0; i < blocks; i++ {
		var notBlock *bitBlock
		block, ok := set.getBlock(i)
		if !ok {
			// create new block and fill it with 1s
			notBlock = newBitBlock()
			notBlock.fillAll()
		} else {
			notBlock = block.notAll()
		}

		result.blocks[i] = notBlock
	}

	//  make sure to And the result with cauldron.AllSet so we don't have dangling 1s
	return result.And(set.cauldron.AllSet())
}
func (set bitSet) AndCount(rhsSet ReadOnlySet) (uint64, error) {
	rhs, ok := rhsSet.(*bitSet)
	if !ok {
		return 0, errors.New("invalid set passed to And")
	}

	set.lock.RLock()
	defer set.lock.RUnlock()

	var count uint64 = 0

	blocks := set.getBlockCount()
	for i := 0; i < blocks; i++ {
		lhsBlock, lhsOk := set.getBlock(i)
		rhsBlock, rhsOk := rhs.getBlock(i)

		if !lhsOk || !rhsOk {
			// dont need to And if either left or right are empty
			continue
		}

		count += lhsBlock.andCount(rhsBlock)
	}

	return count, nil
}

func (set bitSet) OrCount(rhsSet ReadOnlySet) (uint64, error) {
	rhs, ok := rhsSet.(*bitSet)
	if !ok {
		return 0, errors.New("invalid set passed to Or")
	}

	set.lock.RLock()
	defer set.lock.RUnlock()

	var count uint64 = 0

	blocks := set.getBlockCount()
	for i := 0; i < blocks; i++ {
		lhsBlock, lhsOk := set.getBlock(i)
		rhsBlock, rhsOk := rhs.getBlock(i)

		if !lhsOk && !rhsOk {
			// dont need to Or if both left or right are empty
			continue
		} else if !rhsOk {
			count += lhsBlock.countAll()
		} else if !lhsOk {
			count += rhsBlock.countAll()
		} else {
			count += lhsBlock.orCount(rhsBlock)
		}
	}

	return count, nil
}

func (set bitSet) NotCount() (uint64, error) {
	set.lock.RLock()
	defer set.lock.RUnlock()

	result := newBitSet(set.cauldron).(*bitSet)

	blocks := set.getBlockCount()
	for i := 0; i < blocks; i++ {
		var notBlock *bitBlock

		block, ok := set.getBlock(i)
		if !ok {
			// create new block and fill it with 1s
			notBlock = newBitBlock()
			notBlock.fillAll()
		} else {
			notBlock = block.notAll()
		}

		result.blocks[i] = notBlock
	}

	//  make sure to And the result with cauldron.AllSet so we don't have dangling 1s
	cleanNot, err := result.And(set.cauldron.AllSet())
	if err != nil {
		return 0, err
	}

	return cleanNot.Count(), nil
}

func (set bitSet) Count() uint64 {
	set.lock.RLock()
	defer set.lock.RUnlock()

	var count uint64 = 0
	blocks := set.getBlockCount()
	for i := 0; i < blocks; i++ {
		block, ok := set.getBlock(i)
		if ok {
			count += block.countAll()
		}
	}

	return count
}

func (set bitSet) Size() uint64 {
	return set.cauldron.Capacity()
}

func (set bitSet) ToSlice() []Item {
	set.lock.RLock()
	defer set.lock.RUnlock()

	// not the most efficient way to do this,
	//but we'll see if we need to optimize this further later
	var items []Item
	blocks := set.getBlockCount()
	for i := 0; i < blocks; i++ {
		block, ok := set.getBlock(i)
		if ok {
			for j := 0; j < BitsPerBlock; j++ {
				bit, err := block.getBit(j)
				if err == nil && bit {
					index := (uint64(i)*BitsPerBlock) + uint64(j)
					item, err := set.cauldron.ItemFor(index)
					if err == nil {
						items = append(items, item)
					}
				}
			}
		}
	}

	return items
}

func (set bitSet) GetBit(index uint64) (bool, error) {
	set.lock.RLock()
	defer set.lock.RUnlock()

	blockID := int(index / BitsPerBlock)
	idx := int(index % BitsPerBlock)
	block, ok := set.getBlock(blockID)
	if !ok {
		// we don't have that block, so assume all 0s (eg. false the bit is not set)
		return false, nil
	}

	return block.getBit(idx)
}

func (set *bitSet) SetBit(index uint64) error {
	set.lock.Lock()
	defer set.lock.Unlock()

	blockID := int(index / BitsPerBlock)
	idx := int(index % BitsPerBlock)
	block, ok := set.getBlock(blockID)
	if !ok {
		// create the block
		block = newBitBlock()
		set.blocks[blockID] = block
	}

	return block.setBit(idx)
}

func (set *bitSet) UnsetBit(index uint64) error {
	set.lock.Lock()
	defer set.lock.Unlock()

	blockID := int(index / BitsPerBlock)
	idx := int(index % BitsPerBlock)
	block, ok := set.getBlock(blockID)
	if !ok {
		// no need to do anything, it already is assumed to be all 0's
		return nil
	}

	return block.unsetBit(idx)
}

func (set *bitSet) Clear() {
	set.lock.Lock()
	defer set.lock.Unlock()

	blocks := set.getBlockCount()
	for i := 0; i < blocks; i++ {
		block, ok := set.getBlock(i)
		if !ok {
			// no block here, nothing to do
		} else {
			// Note: should we just remove all the blocks (as block=nil assumes all 0's)
			block.clearAll()
		}
	}
}

func (set *bitSet) Fill() {
	set.lock.Lock()
	defer set.lock.Unlock()

	blocks := set.getBlockCount()
	for i := 0; i < blocks; i++ {
		block, ok := set.getBlock(i)
		if !ok {
			// create the block
			block = newBitBlock()
			set.blocks[i] = block
		}
		block.fillAll()
	}

	// the last block will likely have extra dangling 1's
	// options are to AND it with the Cauldron's AllSet or,
	// handle the last block as a special case.
	lastBlock, ok := set.getBlock(blocks-1)
	if ok {
		lastBits := int(set.Size() % BitsPerBlock)
		lastBlock.fill(lastBits)
	}
}

func (set bitSet) getBlockCount() int {
	capacity := set.cauldron.Capacity()
	blocks := int(capacity / BitsPerBlock)
	if capacity%BitsPerBlock > 0 {
		blocks++
	}

	return blocks
}

func (set bitSet) getBlock(id int) (*bitBlock, bool) {
	block, ok := set.blocks[id]

	if !ok || block == nil {
		// The "id" might exist in the map, but the bitBlock == nil
		// Return "false" in this case!!!
		return nil, false
	}

	return block, true
}