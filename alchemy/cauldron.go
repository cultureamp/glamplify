package alchemy

import (
	"sync"

	"github.com/go-errors/errors"
)

// Item type used to represent a record, suggest ID or UUID as a string
type Item string

// Cauldron interface represents the container of all aspect, facets, and items
type Cauldron interface {
	Aspect(name string) (Aspect, error)
	Aspects() ([]Aspect, error)
	NewAspect(name string) (Aspect, error)
	NewAspectWithDisplayName(name string, displayName string) (Aspect, error)

	Capacity() uint64
	Count() uint64

	IndexFor(item Item) (uint64, error)
	ItemFor(index uint64) (Item, error)

	Upsert(item Item) (uint64, error)
	TryRemove(item Item) (bool, error)

	AllSet() ReadOnlySet

	NewSet() Set
}

type bitCauldron struct {
	count    uint64
	capacity uint64

	freeSlots stack
	allSet    Set

	indexToItem map[uint64]Item
	itemToIndex map[Item]uint64

	childAspects map[string]Aspect

	lock *sync.RWMutex
}

func NewBitCauldron(sizeEstimate uint64) Cauldron {
	cauldron := &bitCauldron{
		count:    0,
		capacity: 0,

		freeSlots: newLinkedListStack(),

		indexToItem: make(map[uint64]Item, sizeEstimate),
		itemToIndex: make(map[Item]uint64, sizeEstimate),

		childAspects: make(map[string]Aspect, sizeEstimate),

		lock: &sync.RWMutex{},
	}

	cauldron.allSet = newBitSet(cauldron)

	return cauldron
}

func (cauldron bitCauldron) Aspect(name string) (Aspect, error) {
	cauldron.lock.RLock()
	defer cauldron.lock.RUnlock()

	aspect, ok := cauldron.childAspects[name]
	if !ok {
		return nil, errors.New("no aspect found with that name")
	}

	return aspect, nil
}

func (cauldron bitCauldron) Aspects() ([]Aspect, error) {
	cauldron.lock.RLock()
	defer cauldron.lock.RUnlock()

	len := len(cauldron.childAspects)
	if len == 0 {
		return []Aspect{}, errors.New("no aspects")
	}

	var aspects = make([]Aspect, 0, len)
	for _, aspect := range cauldron.childAspects {
		aspects = append(aspects, aspect)
	}

	return aspects, nil
}

func (cauldron *bitCauldron) NewAspect(name string) (Aspect, error) {
	return cauldron.NewAspectWithDisplayName(name, name)
}

func (cauldron *bitCauldron) NewAspectWithDisplayName(name string, displayName string) (Aspect, error) {
	cauldron.lock.Lock()
	defer cauldron.lock.Unlock()

	aspect, ok := cauldron.childAspects[name]
	if ok {
		return aspect, errors.New("aspect with that name already exists")
	}

	aspect = newBitAspect(name, displayName, cauldron)
	cauldron.childAspects[name] = aspect
	return aspect, nil
}

func (cauldron bitCauldron) Capacity() uint64 {
	return cauldron.capacity
}

func (cauldron bitCauldron) Count() uint64 {
	return cauldron.count
}

func (cauldron bitCauldron) IndexFor(item Item) (uint64, error) {
	cauldron.lock.RLock()
	defer cauldron.lock.RUnlock()

	index, ok := cauldron.itemToIndex[item]
	if !ok {
		return 0, errors.New("no item in lookup")
	}

	return index, nil
}

func (cauldron bitCauldron) ItemFor(index uint64) (Item, error) {
	cauldron.lock.RLock()
	defer cauldron.lock.RUnlock()

	item, ok := cauldron.indexToItem[index]
	if !ok {
		return "", errors.New("no index in lookup")
	}

	return item, nil
}

func (cauldron *bitCauldron) Upsert(item Item) (uint64, error) {

	cauldron.lock.Lock()
	defer cauldron.lock.Unlock()

	index, ok := cauldron.itemToIndex[item]
	if ok {
		// it already exists, so no-op just return it
		return index, nil
	}

	// if isEmpty is a quick test to see if we can re-use an index
	if cauldron.freeSlots.isEmpty() {
		// we don't have a spare slot we can re-use
		index = cauldron.capacity
		cauldron.capacity++
	} else {
		// try and pop a free slot
		var err error
		index, err = cauldron.freeSlots.pop()
		if err != nil {
			// we don't have a spare slot we can re-use
			index = cauldron.capacity
			cauldron.capacity++
		}
	}

	err := cauldron.allSet.SetBit(index)
	cauldron.itemToIndex[item] = index
	cauldron.indexToItem[index] = item

	cauldron.count++
	return index, err
}

func (cauldron *bitCauldron) TryRemove(item Item) (bool, error) {

	index, err := cauldron.IndexFor(item)
	if err != nil {
		// it does not exists, so no-op nothing to do
		return false, nil
	}

	// needs to outside the Lock() as it calls RLock itself
	// and GO doesn't support recursive locks :(
	aspects, _ := cauldron.Aspects()

	cauldron.lock.Lock()
	defer cauldron.lock.Unlock()

	delete(cauldron.itemToIndex, item)
	delete(cauldron.indexToItem, index)

	err = cauldron.allSet.UnsetBit(index)

	// need to set all the bits in this index for all facets to ZERO
	for _, aspect := range aspects{
		facets, _ := aspect.Facets()
		for _, facet := range facets {
			facet.UnsetBitForIndex(index)
		}
	}

	cauldron.freeSlots.push(index)
	cauldron.count--
	return true, err
}

func (cauldron bitCauldron) AllSet() ReadOnlySet {
	return cauldron.allSet
}

func (cauldron *bitCauldron) NewSet() Set {
	return newBitSet(cauldron)
}