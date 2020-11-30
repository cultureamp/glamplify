package alchemy

import (
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"testing"
	"time"

	"gotest.tools/assert"
)

const (
	TestRealWorldLoops = 1000			// 1,000
	TestRealWorldSize = Long(10000000) 	// 10 million
)

func Test_BitCauldron_RealWorld(t *testing.T) {
	caul := newBitCauldron()

	start := time.Now()
	items := make([]Item, 0, TestRealWorldSize)
	for i := Long(0); i < TestRealWorldSize; i++ {
		item := Item(strconv.FormatInt(int64(i), 10))
		caul.Upsert(item)
		items = append(items, item)
	}
	stop := time.Now()
	fmt.Printf("time to upsert %d items = %f sec\n", TestRealWorldSize, stop.Sub(start).Seconds())

	loc, _ := caul.NewAspect("Location")
	vic, _ := loc.NewFacet("Victoria")
	nsw, _ := loc.NewFacet("New South Wales")
	qld, _ := loc.NewFacet("Queensland")
	sa, _ := loc.NewFacet("South Australia")
	wa, _ := loc.NewFacet("Western Australia")
	tas, _ := loc.NewFacet("Tasmania")
	act, _ := loc.NewFacet("Australian Capital Territory")

	start = time.Now()
	for i := Long(0); i < TestRealWorldSize-7; i = i + 7 {
		vic.SetBitForIndex(i)
		nsw.SetBitForIndex(i+1)
		qld.SetBitForIndex(i+2)
		sa.SetBitForIndex(i+3)
		wa.SetBitForIndex(i+4)
		tas.SetBitForIndex(i+5)
		act.SetBitForIndex(i+6)
	}
	stop = time.Now()
	fmt.Printf("time to set up states for %d items = %f sec\n", TestRealWorldSize, stop.Sub(start).Seconds())

	typ, _ := caul.NewAspect("Type")
	mgr, _ := typ.NewFacet("Manager")
	ic, _ := typ.NewFacet("IC")

	start = time.Now()
	for i := Long(0); i < TestRealWorldSize-2; i = i + 2 {
		mgr.SetBitForIndex(i)
		ic.SetBitForIndex(i+1)
	}
	stop = time.Now()
	fmt.Printf("time to set up types for %d items = %f sec\n", TestRealWorldSize, stop.Sub(start).Seconds())

	dept, _ := caul.NewAspect("Department")
	product, _ := dept.NewFacet("Product")
	customer, _ := dept.NewFacet("Customer")
	org, _ := dept.NewFacet("Org")

	start = time.Now()
	for i := Long(0); i < TestRealWorldSize-3; i = i + 3 {
		product.SetBitForIndex(i)
		customer.SetBitForIndex(i+1)
		org.SetBitForIndex(i+2)
	}
	stop = time.Now()
	fmt.Printf("time to set up departments for %d items = %f sec\n", TestRealWorldSize, stop.Sub(start).Seconds())

	var result Set
	var count Long

	start = time.Now()
	for i := 0; i < TestRealWorldLoops; i++ {
		result, _ = vic.Or(nsw)
		result, _ = mgr.AndSet(result)
		count, _ = product.AndCountSet(result)
	}
	stop = time.Now()
	d := stop.Sub(start)
	fmt.Printf("time to vic.Or(nsw).And(mgr).AndCount(product) = %d, %d times for %d items = %f sec (%f ms per op)\n", count, TestRealWorldLoops, TestRealWorldSize, d.Seconds(), float32(d.Milliseconds())/TestRealWorldLoops)
}

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