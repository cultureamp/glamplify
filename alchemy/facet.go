package alchemy

import (
	"github.com/go-errors/errors"
)

type Facet interface {
	Name() string
	DisplayName() string
	Aspect() Aspect
	Set() ReadOnlySet
	Count() uint64
	ToSlice() []Item

	GetBitForItem(item Item) (bool, error)
	GetBitForIndex(index uint64) (bool, error)

	SetBitForItem(item Item) error
	SetBitForIndex(index uint64) error
	UnsetBitForItem(item Item) error
	UnsetBitForIndex(index uint64) error

	And(rhs Facet) (Set, error)
	AndSet(rhs ReadOnlySet) (Set, error)
	Or(rhs Facet) (Set, error)
	OrSet(rhs ReadOnlySet) (Set, error)
	Not() (Set, error)

	AndCount(rhs Facet) (uint64, error)
	AndCountSet(rhs ReadOnlySet) (uint64, error)
	OrCount(rhs Facet) (uint64, error)
	OrCountSet(rhs ReadOnlySet) (uint64, error)
	NotCount() (uint64, error)
}

type bitFacet struct {
	name        string
	displayName string
	aspect      Aspect
	cauldron    Cauldron
	set         Set
}

func newBitFacet(name string, displayName string, aspect Aspect, cauldron Cauldron) Facet {
	return &bitFacet{
		name:        name,
		displayName: displayName,
		aspect:      aspect,
		cauldron:    cauldron,
		set:         cauldron.NewSet(),
	}
}

func (facet bitFacet) Name() string {
	return facet.name
}

func (facet bitFacet) DisplayName() string {
	return facet.displayName
}

func (facet bitFacet) Aspect() Aspect {
	return facet.aspect
}

func (facet bitFacet) Set() ReadOnlySet {
	return facet.set
}

func (facet *bitFacet) Count() uint64 {
	return facet.set.Count()
}

func (facet *bitFacet) ToSlice() []Item {
	return facet.set.ToSlice()
}

func (facet bitFacet) GetBitForItem(item Item) (bool, error) {
	index, err := facet.cauldron.IndexFor(item)
	if err != nil {
		return false, err
	}

	return facet.GetBitForIndex(index)
}

func (facet bitFacet) GetBitForIndex(index uint64) (bool, error) {
	if index >= facet.cauldron.Capacity() {
		return false, errors.New("index greater than cauldron capacity")
	}
	return facet.set.GetBit(index)
}

func (facet *bitFacet) SetBitForItem(item Item) error {
	index, err := facet.cauldron.IndexFor(item)
	if err != nil {
		return err
	}

	return facet.SetBitForIndex(index)
}

func (facet *bitFacet) SetBitForIndex(index uint64) error {
	if index >= facet.cauldron.Capacity() {
		return errors.New("index greater than cauldron capacity")
	}
	return facet.set.SetBit(index)
}

func (facet *bitFacet) UnsetBitForItem(item Item) error {
	index, err := facet.cauldron.IndexFor(item)
	if err != nil {
		return err
	}

	return facet.UnsetBitForIndex(index)
}

func (facet *bitFacet) UnsetBitForIndex(index uint64) error {
	if index >= facet.cauldron.Capacity() {
		return errors.New("index greater than cauldron capacity")
	}
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

func (facet bitFacet) AndCount(rhs Facet) (uint64, error) {
	return facet.set.AndCount(rhs.Set())
}

func (facet bitFacet) AndCountSet(rhs ReadOnlySet) (uint64, error) {
	return facet.set.AndCount(rhs)
}

func (facet bitFacet) OrCount(rhs Facet) (uint64, error) {
	return facet.set.OrCount(rhs.Set())
}

func (facet bitFacet) OrCountSet(rhs ReadOnlySet) (uint64, error) {
	return facet.set.OrCount(rhs)
}

func (facet bitFacet) NotCount() (uint64, error) {
	return facet.set.NotCount()
}
