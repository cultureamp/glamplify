package alchemy

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	TestRealWorldLoops = 1000			// 1,000
	TestRealWorldSize = uint64(10000000) 	// 10 million
)

func Test_BitCauldron_RealWorld(t *testing.T) {
	caul := NewBitCauldron(TestRealWorldSize)

	start := time.Now()
	items := make([]Item, 0, TestRealWorldSize)
	for i := uint64(0); i < TestRealWorldSize; i++ {
		item := Item(strconv.FormatUint(i, 10))
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
	for i := uint64(0); i < TestRealWorldSize-7; i = i + 7 {
		vic.SetBitForIndex(i)
		nsw.SetBitForIndex(i+1)
		qld.SetBitForIndex(i+2)
		sa.SetBitForIndex(i+3)
		wa.SetBitForIndex(i+4)
		tas.SetBitForIndex(i+5)
		act.SetBitForIndex(i+6)
	}
	stop = time.Now()
	fmt.Printf("time to set up states for %s items = %f sec\n", humanize.Comma(int64(TestRealWorldSize)), stop.Sub(start).Seconds())

	role, _ := caul.NewAspect("Role")
	mgr, _ := role.NewFacet("Manager")
	ic, _ := role.NewFacet("IC")

	start = time.Now()
	for i := uint64(0); i < TestRealWorldSize-2; i = i + 2 {
		mgr.SetBitForIndex(i)
		ic.SetBitForIndex(i+1)
	}
	stop = time.Now()
	fmt.Printf("time to set up roles for %s items = %f sec\n", humanize.Comma(int64(TestRealWorldSize)), stop.Sub(start).Seconds())

	dept, _ := caul.NewAspect("Department")
	product, _ := dept.NewFacet("Product")
	customer, _ := dept.NewFacet("Customer")
	org, _ := dept.NewFacet("Org")

	start = time.Now()
	for i := uint64(0); i < TestRealWorldSize-3; i = i + 3 {
		product.SetBitForIndex(i)
		customer.SetBitForIndex(i+1)
		org.SetBitForIndex(i+2)
	}
	stop = time.Now()
	fmt.Printf("time to set up departments for %s items = %f sec\n", humanize.Comma(int64(TestRealWorldSize)), stop.Sub(start).Seconds())

	var result Set
	var count uint64

	start = time.Now()
	for i := 0; i < TestRealWorldLoops; i++ {
		result, _ = vic.Or(nsw)
		result, _ = mgr.AndSet(result)
		count, _ = product.AndCountSet(result)
	}
	stop = time.Now()
	d := stop.Sub(start)
	fmt.Printf("time to vic.Or(nsw).And(mgr).AndCount(product) = %s results called %s times for %s items = %f sec (%f ms per op)\n",  humanize.Comma(int64(count)),  humanize.Comma(TestRealWorldLoops), humanize.Comma(int64(TestRealWorldSize)), d.Seconds(), float32(d.Milliseconds())/TestRealWorldLoops)
}

func Test_New_BitCauldron(t *testing.T) {
	caul := NewBitCauldron(TestRealWorldSize)
	assert.NotNil(t, caul)
}

func Test_BitCauldron_NewAspect(t *testing.T) {
	caul := NewBitCauldron(TestRealWorldSize)

	loc, err := caul.NewAspect("Location")
	assert.Nil(t, err)
	assert.Equal(t, "Location", loc.Name())
	assert.Equal(t, "Location", loc.DisplayName())

	tenure, err := caul.NewAspectWithDisplayName("Ten", "Tenure")
	assert.Nil(t, err)
	assert.Equal(t,  "Ten", tenure.Name())
	assert.Equal(t,  "Tenure", tenure.DisplayName())
}

func Test_BitCauldron_GetAspect(t *testing.T) {
	caul := NewBitCauldron(TestRealWorldSize)

	loc1, err := caul.NewAspect("Location")
	assert.Nil(t, err)

	loc2, err := caul.Aspect("Location")
	assert.Nil(t, err)
	assert.Equal(t, loc2, loc1)

	// doesn't exist
	tenure, err := caul.Aspect("Tenure")
	assert.NotNil(t, err)
	assert.Nil(t, tenure)
}

func Test_BitCauldron_AllAspects(t *testing.T) {
	caul := NewBitCauldron(TestRealWorldSize)

	loc, err := caul.NewAspect("Location")
	assert.Nil(t, err)

	tenure, err := caul.NewAspectWithDisplayName("Ten", "Tenure")
	assert.Nil(t, err)

	aspects, err := caul.Aspects()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(aspects))
	assert.True(t, aspectsContains(aspects, loc))
	assert.True(t, aspectsContains(aspects, tenure))
}

func Test_BitCauldron_TryRemove(t *testing.T) {
	caul := NewBitCauldron(TestRealWorldSize)

	item := Item(uuid.New().String())
	idx, err := caul.Upsert(item)
	assert.Nil(t, err)
	assert.Equal(t, uint64(0), idx)

	loc, err := caul.NewAspect("Location")
	assert.Nil(t, err)

	melb, err := loc.NewFacet("Melbourne")
	assert.Nil(t, err)

	syd, err := loc.NewFacet("Sydney")
	assert.Nil(t, err)

	melb.SetBitForIndex(idx)
	syd.SetBitForIndex(idx)

	ok, err := caul.TryRemove(item)
	assert.Nil(t, err)
	assert.True(t, ok)

	bit, err := melb.GetBitForIndex(idx)
	assert.Nil(t, err)
	assert.False(t, bit)

	bit, err = syd.GetBitForIndex(idx)
	assert.Nil(t, err)
	assert.False(t, bit)

	// remove again should fail
	ok, err = caul.TryRemove(item)
	assert.Nil(t, err)
	assert.False(t, ok)
}

func aspectsContains(s []Aspect, e Aspect) bool {
	for _, f := range s {
		if f == e {
			return true
		}
	}
	return false
}