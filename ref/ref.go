package ref

import (
	"reflect"
	"time"
)

// Bool returns a reference to the value passed as a parameter.
func Bool(i bool) *bool {
	return &i
}

// Int returns a reference to the value passed as a parameter.
func Int(i int) *int {
	return &i
}

// Int64 returns a reference to the value passed as a parameter.
func Int64(i int64) *int64 {
	return &i
}

// String returns a reference to the value passed as a parameter.
func String(i string) *string {
	return &i
}

// Duration returns a reference to the value passed as a parameter.
func Duration(i time.Duration) *time.Duration {
	return &i
}

// Strings returns a new array of pointers that has the same length as the
// supplied array, and whose elements contain pointers to strings in the
// equivalent indexes from the supplied array.
func Strings(ss []string) []*string {
	r := make([]*string, len(ss))
	for i := range ss {
		r[i] = &ss[i]
	}
	return r
}

// ToStructPointer creates an interface{} instance that is a pointer to the
// underlying struct value of the supplied instance of interface{}. So if a{} is
// supplied, this will return &a{}.
func ToStructPointer(a interface{}) interface{} {
	v := reflect.ValueOf(a)
	vp := reflect.New(reflect.TypeOf(a))
	vp.Elem().Set(v)

	return vp.Interface()
}

// ToStruct returns a reference representing the value of the supplied pointer.
// So if *a{} is supplied, a{} will be returned.
func ToStruct(ptr interface{}) interface{} {
	return reflect.ValueOf(ptr).Elem().Interface()
}
