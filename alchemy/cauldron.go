package alchemy

import (
	"github.com/go-errors/errors"
	"sync"
)

type Long uint64
type Item string

type Cauldron interface {
	GetAspect(name string) (Aspect, error)
	GetAspects() ([]Aspect, error)
	NewAspect(name string) (Aspect, error)
	NewAspectWithDisplayName(name string, displayName string) (Aspect, error)

	GetCapacity() Long
	GetCount() Long

	GetIndexFor(item Item) (Long, error)
	GetItemFor(index Long) (Item, error)

	Upsert(item Item) Long
	TryRemove(item Item) bool

	GetEmptySet() ReadOnlySet
	GetAllSet() ReadOnlySet

	NewSet() Set
}

type bitCauldron struct {
	count    Long
	capacity Long

	freeSlots stack
	allSet    Set
	emptySet  Set

	indexToItem map[Long]Item
	itemToIndex map[Item]Long

	childAspects map[string]Aspect

	lock *sync.RWMutex
}

func newBitCauldron() Cauldron {
	cauldron := &bitCauldron{
		count:    0,
		capacity: 0,

		freeSlots: newLinkedListStack(),

		indexToItem: map[Long]Item{},
		itemToIndex: map[Item]Long{},

		childAspects: map[string]Aspect{},

		lock: &sync.RWMutex{},
	}

	cauldron.allSet = newBitSet(cauldron)
	cauldron.emptySet = newBitSet(cauldron)

	return cauldron
}

func (cauldron bitCauldron) GetAspect(name string) (Aspect, error) {
	cauldron.lock.RLock()
	defer cauldron.lock.RUnlock()

	aspect, ok := cauldron.childAspects[name]
	if !ok {
		return nil, errors.New("no aspect found with that name")
	}

	return aspect, nil
}

func (cauldron bitCauldron) GetAspects() ([]Aspect, error) {
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

func (cauldron bitCauldron) GetCapacity() Long {
	cauldron.lock.RLock()
	defer cauldron.lock.RUnlock()

	return cauldron.capacity
}

func (cauldron bitCauldron) GetCount() Long {
	cauldron.lock.RLock()
	defer cauldron.lock.RUnlock()

	return cauldron.count
}

func (cauldron bitCauldron) GetIndexFor(item Item) (Long, error) {
	cauldron.lock.RLock()
	defer cauldron.lock.RUnlock()

	index, ok := cauldron.itemToIndex[item]
	if !ok {
		return 0, errors.New("no item in lookup")
	}

	return index, nil
}

func (cauldron bitCauldron) GetItemFor(index Long) (Item, error) {
	cauldron.lock.RLock()
	defer cauldron.lock.RUnlock()

	item, ok := cauldron.indexToItem[index]
	if !ok {
		return "", errors.New("no index in lookup")
	}

	return item, nil
}

func (cauldron *bitCauldron) Upsert(item Item) Long {

	index, err := cauldron.GetIndexFor(item)
	if err == nil {
		// it already exists, so no-op just return it
		return index
	}

	cauldron.lock.Lock()
	defer cauldron.lock.Unlock()

	// doesn't exist, so add it
	index, err = cauldron.freeSlots.pop()
	if err != nil {
		// we don't have a spare slot we can re-use
		index = cauldron.capacity
		cauldron.capacity++
	}

	cauldron.allSet.SetBit(index)
	cauldron.itemToIndex[item] = index
	cauldron.indexToItem[index] = item

	cauldron.count++
	return index
}

func (cauldron *bitCauldron) TryRemove(item Item) bool {
	cauldron.lock.Lock()
	defer cauldron.lock.Unlock()

	index, err := cauldron.GetIndexFor(item)
	if err != nil {
		// it does not exists, so no-op nothing to do
		return false
	}

	delete(cauldron.itemToIndex, item)
	delete(cauldron.indexToItem, index)

	cauldron.allSet.UnsetBit(index)

	// TODO - can we just AND all results with the cauldron.allSet which will do this for us?
	// need to set all the bits in this index for all facets to ZERO
	aspects, _ := cauldron.GetAspects()
	for _, aspect := range aspects{
		facets, _ := aspect.GetFacets()
		for _, facet := range facets {
			facet.UnseByIndex(index)
		}
	}

	cauldron.freeSlots.push(index)
	cauldron.count--
	return true
}

func (cauldron bitCauldron) GetEmptySet() ReadOnlySet {
	return cauldron.emptySet
}

func (cauldron bitCauldron) GetAllSet() ReadOnlySet {
	return cauldron.allSet
}

func (cauldron *bitCauldron) NewSet() Set {
	return newBitSet(cauldron)
}