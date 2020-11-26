package alchemy

import (
	"sync"
)

type Facet interface {
	GetName() string
	GetDisplayName() string
	GetAspect() Aspect

	Set() ReadOnlySet

	GetByItem(item Item) (bool, error)
	GetByIndex(index Long) (bool, error)

	SetByItem(item Item) error
	SetByIndex(index Long) error
	UnsetByItem(item Item) error
	UnseByIndex(index Long) error

	And(rhs Facet) (Set, error)
	AndSet(rhs ReadOnlySet) (Set, error)
	Or(rhs Facet) (Set, error)
	OrSet(rhs ReadOnlySet) (Set, error)
	Not() (Set, error)

	AndCount(rhs Facet) (Long, error)
	AndCountSet(rhs ReadOnlySet) (Long, error)
	OrCount(rhs Facet) (Long, error)
	OrCountSet(rhs ReadOnlySet) (Long, error)
	NotCount() (Long, error)
}

type bitFacet struct {
	name        string
	displayName string
	aspect      Aspect
	cauldron    *bitCauldron
	set         Set
	lock        *sync.RWMutex
}

func newBitFacet(name string, displayName string, aspect Aspect, cauldron *bitCauldron) Facet {
	return &bitFacet{
		name:        name,
		displayName: displayName,
		aspect:      aspect,
		cauldron:    cauldron,
		set:         cauldron.NewSet(),
		lock:        &sync.RWMutex{},
	}
}

func (facet bitFacet) GetName() string {
	return facet.name
}

func (facet bitFacet) GetDisplayName() string {
	return facet.displayName
}

func (facet bitFacet) GetAspect() Aspect {
	return facet.aspect
}

func (facet bitFacet) Set() ReadOnlySet {
	return facet.set
}

func (facet bitFacet) GetByItem(item Item) (bool, error) {
	index, err := facet.cauldron.GetIndexFor(item)
	if err != nil {
		return false, err
	}

	return facet.GetByIndex(index)
}

func (facet bitFacet) GetByIndex(index Long) (bool, error) {
	return facet.set.GetBit(index)
}

func (facet *bitFacet) SetByItem(item Item) error {
	index, err := facet.cauldron.GetIndexFor(item)
	if err != nil {
		return err
	}

	return facet.SetByIndex(index)
}

func (facet *bitFacet) SetByIndex(index Long) error {
	return facet.set.SetBit(index)
}

func (facet *bitFacet) UnsetByItem(item Item) error {
	index, err := facet.cauldron.GetIndexFor(item)
	if err != nil {
		return err
	}

	return facet.UnseByIndex(index)
}

func (facet *bitFacet) UnseByIndex(index Long) error {
	return facet.set.UnsetBit(index)
}

func (facet bitFacet) And(rhs Facet) (Set, error) {
	return facet.set.And(rhs.Set())
}

func (facet bitFacet) AndSet(rhs ReadOnlySet) (Set, error) {
	return facet.set.And(rhs)
}

func (facet bitFacet) Or(rhs Facet) (Set, error) {
	return facet.set.Or(rhs.Set())
}

func (facet bitFacet) OrSet(rhs ReadOnlySet) (Set, error) {
	return facet.set.Or(rhs)
}

func (facet bitFacet) Not() (Set, error) {
	return facet.set.Not()
}

func (facet bitFacet) AndCount(rhs Facet) (Long, error) {
	return facet.set.AndCount(rhs.Set())
}

func (facet bitFacet) AndCountSet(rhs ReadOnlySet) (Long, error) {
	return facet.set.AndCount(rhs)
}

func (facet bitFacet) OrCount(rhs Facet) (Long, error) {
	return facet.set.OrCount(rhs.Set())
}

func (facet bitFacet) OrCountSet(rhs ReadOnlySet) (Long, error) {
	return facet.set.OrCount(rhs)
}

func (facet bitFacet) NotCount() (Long, error) {
	return facet.set.NotCount()
}
