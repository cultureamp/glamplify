package alchemy

import (
	"testing"

	"gotest.tools/assert"
)

func Test_New_BitAspect(t *testing.T) {
	aspect := newBitAspect("Location", "Location", testCauldron)
	assert.Assert(t, aspect != nil, aspect)
}

func Test_BitAspect_Name_DisplayName(t *testing.T) {
	aspect := newBitAspect("Loc", "Location", testCauldron)
	assert.Assert(t, aspect.Name() == "Loc", aspect)
	assert.Assert(t, aspect.DisplayName() == "Location", aspect)
}

func Test_BitAspect_NewFacet(t *testing.T) {
	aspect := newBitAspect("Loc", "Location", testCauldron)

	melb, err := aspect.NewFacet("Melbourne")
	assert.Assert(t, err == nil, err)
	assert.Assert(t, melb.Name() == "Melbourne", melb.Name())
	assert.Assert(t, melb.DisplayName() == "Melbourne", melb.DisplayName())

	syd, err := aspect.NewFacetWithDisplayName("Syd", "Sydney")
	assert.Assert(t, err == nil, err)
	assert.Assert(t, syd.Name() == "Syd", syd.Name())
	assert.Assert(t, syd.DisplayName() == "Sydney", syd.DisplayName())
}

func Test_BitAspect_NewFacet_Error(t *testing.T) {
	aspect := newBitAspect("Loc", "Location", testCauldron)

	melb1, err := aspect.NewFacet("Melbourne")
	assert.Assert(t, err == nil, err)

	// Already exists, return same facet + error
	melb2, err := aspect.NewFacet("Melbourne")
	assert.Assert(t, err != nil, err)

	assert.Assert(t, melb1 == melb2, melb2)
}

func Test_BitAspect_GetFacet(t *testing.T) {
	aspect := newBitAspect("Loc", "Location", testCauldron)

	melb1, err := aspect.NewFacet("Melbourne")
	assert.Assert(t, err == nil, err)

	melb2, err := aspect.Facet("Melbourne")
	assert.Assert(t, err == nil, err)
	assert.Assert(t, melb1 == melb2, melb2)

	// doesn't exist
	syd, err := aspect.Facet("Sydney")
	assert.Assert(t, err != nil, err)
	assert.Assert(t, syd == nil, err)
}

func Test_BitAspect_AllFacets(t *testing.T) {
	aspect := newBitAspect("Loc", "Location", testCauldron)

	melb, err := aspect.NewFacet("Melbourne")
	assert.Assert(t, err == nil, err)

	syd, err := aspect.NewFacet("Sydney")
	assert.Assert(t, err == nil, err)

	facets, err := aspect.Facets()
	assert.Assert(t, err == nil, err)
	assert.Assert(t, len(facets) == 2, len(facets))
	assert.Assert(t, facetsContains(facets, melb))
	assert.Assert(t, facetsContains(facets, syd))
}

func facetsContains(s []Facet, e Facet) bool {
	for _, f := range s {
		if f == e {
			return true
		}
	}
	return false
}