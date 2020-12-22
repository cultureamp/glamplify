package alchemy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New_BitAspect(t *testing.T) {
	aspect := newBitAspect("Location", "Location", testCauldron)
	assert.NotNil(t, aspect)
}

func Test_BitAspect_Name_DisplayName(t *testing.T) {
	aspect := newBitAspect("Loc", "Location", testCauldron)
	assert.Equal(t, "Loc", aspect.Name())
	assert.Equal(t, "Location", aspect.DisplayName())
}

func Test_BitAspect_NewFacet(t *testing.T) {
	aspect := newBitAspect("Loc", "Location", testCauldron)

	melb, err := aspect.NewFacet("Melbourne")
	assert.Nil(t, err)
	assert.Equal(t, "Melbourne", melb.Name())
	assert.Equal(t, "Melbourne", melb.DisplayName())

	syd, err := aspect.NewFacetWithDisplayName("Syd", "Sydney")
	assert.Nil(t, err)
	assert.Equal(t, "Syd", syd.Name())
	assert.Equal(t, "Sydney", syd.DisplayName())
}

func Test_BitAspect_NewFacet_Error(t *testing.T) {
	aspect := newBitAspect("Loc", "Location", testCauldron)

	melb1, err := aspect.NewFacet("Melbourne")
	assert.Nil(t, err)

	// Already exists, return same facet + error
	melb2, err := aspect.NewFacet("Melbourne")
	assert.NotNil(t, err)
	assert.Equalf(t, melb1, melb2, "Should be equal but got '%s' != '%s'", melb1, melb2)
}

func Test_BitAspect_GetFacet(t *testing.T) {
	aspect := newBitAspect("Loc", "Location", testCauldron)

	melb1, err := aspect.NewFacet("Melbourne")
	assert.Nil(t, err)

	melb2, err := aspect.Facet("Melbourne")
	assert.Nil(t, err)
	assert.Equalf(t, melb1, melb2, "Should be equal but got '%s' != '%s'", melb1, melb2)

	// doesn't exist
	syd, err := aspect.Facet("Sydney")
	assert.NotNil(t, err)
	assert.Nil(t, syd)
}

func Test_BitAspect_AllFacets(t *testing.T) {
	aspect := newBitAspect("Loc", "Location", testCauldron)

	melb, err := aspect.NewFacet("Melbourne")
	assert.Nil(t, err)

	syd, err := aspect.NewFacet("Sydney")
	assert.Nil(t, err)

	facets, err := aspect.Facets()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(facets))
	assert.True(t, facetsContains(facets, melb))
	assert.True(t, facetsContains(facets, syd))
}

func facetsContains(s []Facet, e Facet) bool {
	for _, f := range s {
		if f == e {
			return true
		}
	}
	return false
}
