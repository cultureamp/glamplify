package alchemy

import (
	"github.com/go-errors/errors"
)

type Facet interface {
	Name() string
	DisplayName() string
	Aspect() Aspect
	Set() ReadOnlySet
	Count() Long
	ToSlice() []Item

	GetBitForItem(item Item) (bool, error)
	GetBitForIndex(index Long) (bool, error)

	SetBitForItem(item Item) error
	SetBitForIndex(index Long) error
	UnsetBitForItem(item Item) error
	UnsetBitForIndex(index Long) error

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

func (facet *bitFacet) Count() Long {
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

func (facet bitFacet) GetBitForIndex(index Long) (bool, error) {
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

func (facet *bitFacet) SetBitForIndex(index Long) error {
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

func (facet *bitFacet) UnsetBitForIndex(index Long) error {
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
