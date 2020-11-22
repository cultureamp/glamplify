package alchemy


type Aspect interface {
	GetName() string
	GetDisplayName() string

	NewFacet(name string) Facet
	NewFacetWithDisplayName(name string, displayName string) Facet
	GetFacet(name string) Facet
	GetFacets() []Facet
}

