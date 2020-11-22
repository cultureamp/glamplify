package alchemy

type Long uint64

type Cauldron interface {
	GetAspect(name string) Aspect
	GetAspects() []Aspect
	NewAspect(name string) Aspect
	NewAspectWithDisplayName(name string, displayName string) Aspect

	GetCapacity() Long
	GetCount() Long

	GetIndexFor(item Item) Long
	GetItemFor(index Long) Item

	Upsert(item Item) Long
	TryRemove(item Item) bool

	GetEmptySet() ReadOnlySet
	GetAllSet() ReadOnlySet
}

