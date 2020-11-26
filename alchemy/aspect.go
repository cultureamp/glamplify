package alchemy

import (
	"sync"

	"github.com/go-errors/errors"
)

type Aspect interface {
	Name() string
	DisplayName() string

	NewFacet(name string) (Facet, error)
	NewFacetWithDisplayName(name string, displayName string) (Facet, error)

	Facet(name string) (Facet, error)
	Facets() ([]Facet, error)
}

type bitAspect struct {
	name        string
	displayName string
	cauldron    Cauldron

	childFacets map[string]Facet

	lock *sync.RWMutex
}

func newBitAspect(name string, displayName string, cauldron Cauldron) Aspect {
	return &bitAspect{
		name:        name,
		displayName: displayName,
		cauldron:    cauldron,
		childFacets: map[string]Facet{},
		lock:        &sync.RWMutex{},
	}
}

func (aspect bitAspect) Name() string {
	return aspect.name
}

func (aspect bitAspect) DisplayName() string {
	return aspect.displayName
}

func (aspect bitAspect) NewFacet(name string) (Facet, error) {
	return aspect.NewFacetWithDisplayName(name, name)
}

func (aspect bitAspect) NewFacetWithDisplayName(name string, displayName string) (Facet, error) {
	aspect.lock.Lock()
	defer aspect.lock.Unlock()

	facet, ok := aspect.childFacets[name]
	if ok {
		return facet, errors.New("facet with that name already exists")
	}

	facet = newBitFacet(name, displayName, aspect, aspect.cauldron)
	aspect.childFacets[name] = facet
	return facet, nil
}

func (aspect bitAspect) Facet(name string) (Facet, error) {
	aspect.lock.RLock()
	defer aspect.lock.RUnlock()

	facet, ok := aspect.childFacets[name]
	if !ok {
		return nil, errors.New("no facet found with that name")
	}

	return facet, nil
}

func (aspect bitAspect) Facets() ([]Facet, error) {
	aspect.lock.RLock()
	defer aspect.lock.RUnlock()

	len := len(aspect.childFacets)
	if len == 0 {
		return []Facet{}, errors.New("no facets for aspect")
	}

	var facets = make([]Facet, 0, len)
	for _, facet := range aspect.childFacets {
		facets = append(facets, facet)
	}

	return facets, nil
}
