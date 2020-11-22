package alchemy

type ReadOnlySet interface {
	and(set ReadOnlySet) Set
	or(set ReadOnlySet) Set
	not() Set

	getCount() Long
	getSize() Long

	getBit(index Long) bool

	andCount(set ReadOnlySet) Long
	orCount(set ReadOnlySet) Long
	notCount() Long

	toSlice() []Item

}

type Set interface {
	and(set ReadOnlySet) Set
	or(set ReadOnlySet) Set
	not() Set

	getCount() Long
	getSize() Long

	getBit(index Long) bool

	andCount(set ReadOnlySet) Long
	orCount(set ReadOnlySet) Long
	notCount() Long

	toSlice() []Item

	setBit(index Long)
	unsetBit(index Long)
	clear()
	fill()
}
