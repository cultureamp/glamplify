package alchemy

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_New_BitFacet(t *testing.T) {
	aspect := newBitAspect("Location", "Location", testCauldron)
	facet := newBitFacet("Melbourne", "Melbourne", aspect, testCauldron)
	assert.NotNil(t, facet)
}

func Test_BitFacet_Getters(t *testing.T) {
	aspect := newBitAspect("Location", "Location", testCauldron)
	facet := newBitFacet("Melbourne", "Melbourne", aspect, testCauldron)

	assert.Equal(t, "Melbourne", facet.Name())
	assert.Equal(t, "Melbourne", facet.DisplayName())

	a := facet.Aspect()
	assert.NotNil(t, a)
	assert.Equal(t, aspect.Name(), a.Name())

	s := facet.Set()
	assert.NotNil(t, s)
	assert.Equal(t, uint64(0), s.Count())
}

func Test_BitFacet_Set_Get_Unset_ByIndex(t *testing.T) {

	item := Item(uuid.New().String())
	idx, err := testCauldron.Upsert(item)
	assert.Nil(t, err)

	aspect := newBitAspect("Location", "Location", testCauldron)
	facet := newBitFacet("Melbourne", "Melbourne", aspect, testCauldron)

	bit, err := facet.GetBitForIndex(idx)
	assert.Nil(t, err)
	assert.False(t, bit)

	err = facet.SetBitForIndex(idx)
	assert.Nil(t, err)

	bit, err = facet.GetBitForIndex(idx)
	assert.Nil(t, err)
	assert.True(t, bit)

	err = facet.UnsetBitForIndex(idx)
	assert.Nil(t, err)

	bit, err = facet.GetBitForIndex(idx)
	assert.Nil(t, err)
	assert.False(t, bit)

	// out of bounds
	bit, err = facet.GetBitForIndex(TestSetSize+1)
	assert.NotNil(t, err)
	assert.False(t, bit)

	err = facet.SetBitForIndex(TestSetSize+1)
	assert.NotNil(t, err)

	err = facet.UnsetBitForIndex(TestSetSize+1)
	assert.NotNil(t, err)
}

func Test_BitFacet_Set_Get_Unset_ByItem(t *testing.T) {

	item := Item(uuid.New().String())
	testCauldron.Upsert(item)

	aspect := newBitAspect("Location", "Location", testCauldron)
	facet := newBitFacet("Melbourne", "Melbourne", aspect, testCauldron)

	bit, err := facet.GetBitForItem(item)
	assert.Nil(t, err)
	assert.False(t, bit)

	err = facet.SetBitForItem(item)
	assert.Nil(t, err)

	bit, err = facet.GetBitForItem(item)
	assert.Nil(t, err)
	assert.True(t, bit)

	err = facet.UnsetBitForItem(item)
	assert.Nil(t, err)

	bit, err = facet.GetBitForItem(item)
	assert.Nil(t, err)
	assert.False(t, bit)

	item = Item(uuid.New().String())
	bit, err = facet.GetBitForItem(item)
	assert.NotNil(t, err)
	assert.False(t, bit)

	err = facet.SetBitForItem(item)
	assert.NotNil(t, err)

	err = facet.UnsetBitForItem(item)
	assert.NotNil(t, err)
}

func Test_BitFacet_ToSlice(t *testing.T) {

	location := newBitAspect("Location", "Location", testCauldron)
	melbourne := newBitFacet("Melbourne", "Melbourne", location, testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		idx := uint64(rand.Int63n(TestSetSize))
		melbourne.SetBitForIndex(idx)
	}

	slice := melbourne.ToSlice()
	assert.Equal(t, melbourne.Count(), uint64(len(slice)))
}

func Test_BitFacet_And(t *testing.T) {

	location := newBitAspect("Location", "Location", testCauldron)
	melbourne := newBitFacet("Melbourne", "Melbourne", location, testCauldron)

	tenure := newBitAspect("Tenure", "Tenure", testCauldron)
	oneYear := newBitFacet("OneYear", "One Year", tenure, testCauldron)

	// And two empty sets
	result, err := melbourne.And(oneYear)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint64(0), result.Count())

	// One empty, the other with values
	empty := newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		idx := uint64(rand.Int63n(TestSetSize))
		melbourne.SetBitForIndex(idx)
		oneYear.SetBitForIndex(idx)
	}

	result, err = melbourne.AndSet(empty)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint64(0), result.Count())

	result, err = oneYear.AndSet(empty)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint64(0), result.Count())

	result, err = melbourne.And(oneYear)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, melbourne.Count(), result.Count())
	assert.Equal(t, oneYear.Count(),result.Count())
}

func Test_BitFacet_AndCount(t *testing.T) {

	location := newBitAspect("Location", "Location", testCauldron)
	melbourne := newBitFacet("Melbourne", "Melbourne", location, testCauldron)

	tenure := newBitAspect("Tenure", "Tenure", testCauldron)
	oneYear := newBitFacet("OneYear", "One Year", tenure, testCauldron)

	// And two empty sets
	count, err := melbourne.AndCount(oneYear)
	assert.Nil(t, err)
	assert.Equal(t, uint64(0), count)

	// One empty, the other with values
	empty := newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		idx := uint64(rand.Int63n(TestSetSize))
		melbourne.SetBitForIndex(idx)
		oneYear.SetBitForIndex(idx)
	}

	count, err = melbourne.AndCountSet(empty)
	assert.Nil(t, err)
	assert.Equal(t, uint64(0), count)

	count, err = oneYear.AndCountSet(empty)
	assert.Nil(t, err)
	assert.Equal(t, uint64(0), count)

	count, err = melbourne.AndCount(oneYear)
	assert.Nil(t, err)
	assert.Equal(t, count, melbourne.Count())
	assert.Equal(t, count, oneYear.Count())
}

func Test_BitFacet_Or(t *testing.T) {
	location := newBitAspect("Location", "Location", testCauldron)
	melbourne := newBitFacet("Melbourne", "Melbourne", location, testCauldron)

	tenure := newBitAspect("Tenure", "Tenure", testCauldron)
	oneYear := newBitFacet("OneYear", "One Year", tenure, testCauldron)

	// Or two empty sets
	result, err := oneYear.Or(melbourne)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint64(0), result.Count())

	// One empty, the other with values
	empty := newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		idx := uint64(rand.Int63n(TestSetSize))
		melbourne.SetBitForIndex(idx)
		oneYear.SetBitForIndex(idx)
	}

	result, err = melbourne.OrSet(empty)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.Count(), melbourne.Count())

	result, err = melbourne.Or(oneYear)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, melbourne.Count(), result.Count())
	assert.Equal(t, oneYear.Count(), result.Count())
}

func Test_BitFacet_OrCount(t *testing.T) {
	location := newBitAspect("Location", "Location", testCauldron)
	melbourne := newBitFacet("Melbourne", "Melbourne", location, testCauldron)

	tenure := newBitAspect("Tenure", "Tenure", testCauldron)
	oneYear := newBitFacet("OneYear", "One Year", tenure, testCauldron)

	// Or two empty sets
	count, err := oneYear.OrCount(melbourne)
	assert.Nil(t, err)
	assert.Equal(t, uint64(0), count)

	// One empty, the other with values
	empty := newBitSet(testCauldron)

	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		idx := uint64(rand.Int63n(TestSetSize))
		melbourne.SetBitForIndex(idx)
		oneYear.SetBitForIndex(idx)
	}

	count, err = melbourne.OrCountSet(empty)
	assert.Nil(t, err)
	assert.Equal(t, count, melbourne.Count())

	count, err = melbourne.OrCount(oneYear)
	assert.Nil(t, err)
	assert.Equal(t, count, melbourne.Count())
	assert.Equal(t, count, oneYear.Count())
}

func Test_BitFacet_Not(t *testing.T) {

	location := newBitAspect("Location", "Location", testCauldron)
	melbourne := newBitFacet("Melbourne", "Melbourne", location, testCauldron)

	cauldronCount := testCauldron.Count()

	// Not empty set
	result, err := melbourne.Not()
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cauldronCount, result.Count())

	// Add some values
	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := uint64(rand.Int63n(TestSetSize))
		melbourne.SetBitForIndex(bitIdx)
	}

	countBeforeNot := melbourne.Count()

	result, err = melbourne.Not()
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cauldronCount-countBeforeNot, result.Count())
}

func Test_BitFacet_NotCount(t *testing.T) {

	location := newBitAspect("Location", "Location", testCauldron)
	melbourne := newBitFacet("Melbourne", "Melbourne", location, testCauldron)

	cauldronCount := testCauldron.Count()

	// Not empty set
	count, err := melbourne.NotCount()
	assert.Nil(t, err)
	assert.Equal(t, cauldronCount, count)

	// Add some values
	s1 := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(s1)

	for i := 0; i < TestNumberOfBits; i++ {
		bitIdx := uint64(rand.Int63n(TestSetSize))
		melbourne.SetBitForIndex(bitIdx)
	}

	countBeforeNot := melbourne.Count()

	count, err = melbourne.NotCount()
	assert.Nil(t, err)
	assert.Equal(t, cauldronCount-countBeforeNot, count)
}