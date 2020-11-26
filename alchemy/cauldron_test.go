package alchemy

import (
	"github.com/google/uuid"
	"testing"

	"gotest.tools/assert"
)

func Test_New_BitCauldron(t *testing.T) {
	caul := newBitCauldron()
	assert.Assert(t, caul != nil, caul)
}

func Test_BitCauldron_NewAspect(t *testing.T) {
	caul := newBitCauldron()

	loc, err := caul.NewAspect("Location")
	assert.Assert(t, err == nil, err)
	assert.Assert(t, loc.Name() == "Location", loc.Name())
	assert.Assert(t, loc.DisplayName() == "Location", loc.DisplayName())

	tenure, err := caul.NewAspectWithDisplayName("Ten", "Tenure")
	assert.Assert(t, err == nil, err)
	assert.Assert(t, tenure.Name() == "Ten", tenure.Name())
	assert.Assert(t, tenure.DisplayName() == "Tenure", tenure.DisplayName())
}

func Test_BitCauldron_GetAspect(t *testing.T) {
	caul := newBitCauldron()

	loc1, err := caul.NewAspect("Location")
	assert.Assert(t, err == nil, err)

	loc2, err := caul.Aspect("Location")
	assert.Assert(t, err == nil, err)
	assert.Assert(t, loc1 == loc2, loc2)

	// doesn't exist
	tenure, err := caul.Aspect("Tenure")
	assert.Assert(t, err != nil, err)
	assert.Assert(t, tenure == nil, err)
}

func Test_BitCauldron_AllAspects(t *testing.T) {
	caul := newBitCauldron()

	loc, err := caul.NewAspect("Location")
	assert.Assert(t, err == nil, err)

	tenure, err := caul.NewAspectWithDisplayName("Ten", "Tenure")
	assert.Assert(t, err == nil, err)

	aspects, err := caul.Aspects()
	assert.Assert(t, err == nil, err)
	assert.Assert(t, len(aspects) == 2, len(aspects))
	assert.Assert(t, aspectsContains(aspects, loc))
	assert.Assert(t, aspectsContains(aspects, tenure))
}

func Test_BitCauldron_TryRemove(t *testing.T) {
	caul := newBitCauldron()

	item := Item(uuid.New().String())
	idx := caul.Upsert(item)
	assert.Assert(t, idx == 0, idx)

	loc, err := caul.NewAspect("Location")
	assert.Assert(t, err == nil, err)

	melb, err := loc.NewFacet("Melbourne")
	assert.Assert(t, err == nil, err)

	syd, err := loc.NewFacet("Sydney")
	assert.Assert(t, err == nil, err)

	melb.SetBitForIndex(idx)
	syd.SetBitForIndex(idx)

	ok := caul.TryRemove(item)
	assert.Assert(t, ok, ok)

	bit, err := melb.GetBitForIndex(idx)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)

	bit, err = syd.GetBitForIndex(idx)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)

	// remove again should fail
	ok = caul.TryRemove(item)
	assert.Assert(t, !ok, ok)
}

func aspectsContains(s []Aspect, e Aspect) bool {
	for _, f := range s {
		if f == e {
			return true
		}
	}
	return false
}