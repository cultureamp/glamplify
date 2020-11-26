package alchemy

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"gotest.tools/assert"
)

func Test_New_BitFacet(t *testing.T) {
	aspect := newBitAspect("Location", "Location", testCauldron)
	facet := newBitFacet("Melbourne", "Melbourne", aspect, testCauldron)
	assert.Assert(t, facet != nil, facet)
}

func Test_BitFacet_Getters(t *testing.T) {
	aspect := newBitAspect("Location", "Location", testCauldron)
	facet := newBitFacet("Melbourne", "Melbourne", aspect, testCauldron)

	assert.Assert(t, facet.Name() == "Melbourne", facet.Name())
	assert.Assert(t, facet.DisplayName() == "Melbourne", facet.DisplayName())

	a := facet.Aspect()
	assert.Assert(t, a != nil, a)
	assert.Assert(t, a.Name() == aspect.Name(), a.Name())

	s := facet.Set()
	assert.Assert(t, s != nil, s)
	assert.Assert(t, s.Count() == 0, s.Count())
}

func Test_BitFacet_Set_Get_Unset_ByIndex(t *testing.T) {

	item := Item(uuid.New().String())
	idx := testCauldron.Upsert(item)

	aspect := newBitAspect("Location", "Location", testCauldron)
	facet := newBitFacet("Melbourne", "Melbourne", aspect, testCauldron)

	bit, err := facet.GetBitForIndex(idx)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)

	err = facet.SetBitForIndex(idx)
	assert.Assert(t, err == nil, err)

	bit, err = facet.GetBitForIndex(idx)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)

	err = facet.UnsetBitForIndex(idx)
	assert.Assert(t, err == nil, err)

	bit, err = facet.GetBitForIndex(idx)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)

	// out of bounds
	bit, err = facet.GetBitForIndex(TestSetMaxSize+1)
	assert.Assert(t, err != nil, err)
	assert.Assert(t, !bit, bit)

	err = facet.SetBitForIndex(TestSetMaxSize+1)
	assert.Assert(t, err != nil, err)

	err = facet.UnsetBitForIndex(TestSetMaxSize+1)
	assert.Assert(t, err != nil, err)
}

func Test_BitFacet_Set_Get_Unset_ByItem(t *testing.T) {

	item := Item(uuid.New().String())
	testCauldron.Upsert(item)

	aspect := newBitAspect("Location", "Location", testCauldron)
	facet := newBitFacet("Melbourne", "Melbourne", aspect, testCauldron)

	bit, err := facet.GetBitForItem(item)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)

	err = facet.SetBitForItem(item)
	assert.Assert(t, err == nil, err)

	bit, err = facet.GetBitForItem(item)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, bit, bit)

	err = facet.UnsetBitForItem(item)
	assert.Assert(t, err == nil, err)

	bit, err = facet.GetBitForItem(item)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, !bit, bit)

	item = Item(uuid.New().String())
	bit, err = facet.GetBitForItem(item)
	assert.Assert(t, err != nil, err)
	assert.Assert(t, !bit, bit)

	err = facet.SetBitForItem(item)
	assert.Assert(t, err != nil, err)

	err = facet.UnsetBitForItem(item)
	assert.Assert(t, err != nil, err)
}

func Test_BitFacet_ToSlice(t *testing.T) {

	location := newBitAspect("Location", "Location", testCauldron)
	melbourne := newBitFacet("Melbourne", "Melbourne", location, testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		idx := Long(rand.Int63n(TestSetMaxSize))
		melbourne.SetBitForIndex(idx)
	}

	slice := melbourne.ToSlice()
	assert.Assert(t, Long(len(slice)) == melbourne.Count(), len(slice))
}

func Test_BitFacet_And(t *testing.T) {

	location := newBitAspect("Location", "Location", testCauldron)
	melbourne := newBitFacet("Melbourne", "Melbourne", location, testCauldron)

	tenure := newBitAspect("Tenure", "Tenure", testCauldron)
	oneYear := newBitFacet("OneYear", "One Year", tenure, testCauldron)

	// And two empty sets
	result, err := melbourne.And(oneYear)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, result != nil, result)
	assert.Assert(t, result.Count() == 0, result.Count())

	// One empty, the other with values
	empty := newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		idx := Long(rand.Int63n(TestSetMaxSize))
		melbourne.SetBitForIndex(idx)
		oneYear.SetBitForIndex(idx)
	}

	result, err = melbourne.AndSet(empty)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, result != nil, result)
	assert.Assert(t, result.Count() == 0, result.Count())

	result, err = oneYear.AndSet(empty)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, result != nil, result)
	assert.Assert(t, result.Count() == 0, result.Count())

	result, err = melbourne.And(oneYear)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, result != nil, result)
	assert.Assert(t, result.Count() == melbourne.Count(), result.Count())
	assert.Assert(t, result.Count() == oneYear.Count(), result.Count())
}

func Test_BitFacet_AndCount(t *testing.T) {

	location := newBitAspect("Location", "Location", testCauldron)
	melbourne := newBitFacet("Melbourne", "Melbourne", location, testCauldron)

	tenure := newBitAspect("Tenure", "Tenure", testCauldron)
	oneYear := newBitFacet("OneYear", "One Year", tenure, testCauldron)

	// And two empty sets
	count, err := melbourne.AndCount(oneYear)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, count == 0, count)

	// One empty, the other with values
	empty := newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		idx := Long(rand.Int63n(TestSetMaxSize))
		melbourne.SetBitForIndex(idx)
		oneYear.SetBitForIndex(idx)
	}

	count, err = melbourne.AndCountSet(empty)
	assert.Assert(t, err == nil, err)
	assert.Assert(t,count == 0, count)

	count, err = oneYear.AndCountSet(empty)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, count == 0, count)

	count, err = melbourne.AndCount(oneYear)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, count == melbourne.Count(), count)
	assert.Assert(t, count == oneYear.Count(), count)
}

func Test_BitFacet_Or(t *testing.T) {
	location := newBitAspect("Location", "Location", testCauldron)
	melbourne := newBitFacet("Melbourne", "Melbourne", location, testCauldron)

	tenure := newBitAspect("Tenure", "Tenure", testCauldron)
	oneYear := newBitFacet("OneYear", "One Year", tenure, testCauldron)

	// Or two empty sets
	result, err := oneYear.Or(melbourne)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, result != nil, result)
	assert.Assert(t, result.Count() == 0, result.Count())

	// One empty, the other with values
	empty := newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		idx := Long(rand.Int63n(TestSetMaxSize))
		melbourne.SetBitForIndex(idx)
		oneYear.SetBitForIndex(idx)
	}

	result, err = melbourne.OrSet(empty)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, result != nil, result)
	assert.Assert(t, result.Count() == melbourne.Count(), result.Count())

	result, err = melbourne.Or(oneYear)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, result != nil, result)
	assert.Assert(t, result.Count() == melbourne.Count(), result.Count())
	assert.Assert(t, result.Count() == oneYear.Count(), result.Count())
}

func Test_BitFacet_OrCount(t *testing.T) {
	location := newBitAspect("Location", "Location", testCauldron)
	melbourne := newBitFacet("Melbourne", "Melbourne", location, testCauldron)

	tenure := newBitAspect("Tenure", "Tenure", testCauldron)
	oneYear := newBitFacet("OneYear", "One Year", tenure, testCauldron)

	// Or two empty sets
	count, err := oneYear.OrCount(melbourne)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, count == 0, count)

	// One empty, the other with values
	empty := newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		idx := Long(rand.Int63n(TestSetMaxSize))
		melbourne.SetBitForIndex(idx)
		oneYear.SetBitForIndex(idx)
	}

	count, err = melbourne.OrCountSet(empty)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, count == melbourne.Count(), count)

	count, err = melbourne.OrCount(oneYear)
	assert.Assert(t, err == nil, err)
	assert.Assert(t, count == melbourne.Count(), count)
	assert.Assert(t, count == oneYear.Count(), count)
}

func Test_BitFacet_Not(t *testing.T) {

	location := newBitAspect("Location", "Location", testCauldron)
	melbourne := newBitFacet("Melbourne", "Melbourne", location, testCauldron)

	cauldronCount := testCauldron.Count()

	// Not empty set
	result, err := melbourne.Not()
	assert.Assert(t, err == nil, err)
	assert.Assert(t, result != nil, result)
	assert.Assert(t, result.Count() == cauldronCount, result.Count())

	// Add some values
	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := Long(rand.Int63n(TestSetMaxSize))
		melbourne.SetBitForIndex(bitIdx)
	}

	countBeforeNot := melbourne.Count()

	result, err = melbourne.Not()
	assert.Assert(t, err == nil, err)
	assert.Assert(t, result != nil, result)
	assert.Assert(t, result.Count() == (cauldronCount-countBeforeNot), result.Count())
}

func Test_BitFacet_NotCount(t *testing.T) {

	location := newBitAspect("Location", "Location", testCauldron)
	melbourne := newBitFacet("Melbourne", "Melbourne", location, testCauldron)

	cauldronCount := testCauldron.Count()

	// Not empty set
	count, err := melbourne.NotCount()
	assert.Assert(t, err == nil, err)
	assert.Assert(t, count == cauldronCount, count)

	// Add some values
	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := Long(rand.Int63n(TestSetMaxSize))
		melbourne.SetBitForIndex(bitIdx)
	}

	countBeforeNot := melbourne.Count()

	count, err = melbourne.NotCount()
	assert.Assert(t, err == nil, err)
	assert.Assert(t, count == (cauldronCount-countBeforeNot), count)
}