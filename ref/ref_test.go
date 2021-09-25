package ref_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cultureamp/glamplify/ref"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ExampleType struct {
	Member *string
}

func ExampleString() {
	val := "value1"

	ex1 := ExampleType{
		Member: &val,
	}

	ex2 := ExampleType{
		Member: ref.String("value2"),
	}

	fmt.Printf("%s %s\n", *ex1.Member, *ex2.Member)
	// Output: value1 value2
}

func Test_Int(t *testing.T) {
	i := int(49392)
	ptr := ref.Int(i)

	require.NotNil(t, ptr)
	assert.Equal(t, i, *ptr)
}

func Test_Int64(t *testing.T) {
	i := int64(49392215732365141)
	ptr := ref.Int64(i)

	require.NotNil(t, ptr)
	assert.Equal(t, i, *ptr)
}

func Test_Bool(t *testing.T) {
	b := true
	ptr := ref.Bool(b)

	require.NotNil(t, ptr)
	assert.Equal(t, b, *ptr)
}

func Test_Duration(t *testing.T) {
	d := time.Millisecond
	ptr := ref.Duration(d)

	require.NotNil(t, ptr)
	assert.Equal(t, d, *ptr)
}

func Test_String(t *testing.T) {
	str := "flamingo"
	strPtr := ref.String(str)

	require.NotNil(t, strPtr)
	assert.Equal(t, str, *strPtr)
}

func Test_Strings(t *testing.T) {
	slice := []string{"flamingo", "swan", "shag"}
	slicePtr := ref.Strings(slice)

	require.NotNil(t, slicePtr)
	require.Len(t, slicePtr, len(slice))
	actual := make([]string, len(slicePtr))
	for i, e := range slicePtr {
		actual[i] = *e
	}
	assert.Equal(t, slice, actual)
}

type mirror struct {
	who      string
	greatest bool
}

func Test_StructPtr(t *testing.T) {
	start := interface{}(mirror{who: "goat", greatest: true})
	finish := ref.ToStruct(ref.ToStructPointer(start))

	assert.Equal(t, start, finish)
}
