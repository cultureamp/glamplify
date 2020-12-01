package alchemy

import (
	"sync"

	"github.com/go-errors/errors"
)

type node struct {
	data uint64
	next *node
}

type stack interface {
	isEmpty() bool
	push(index uint64)
	pop() (uint64, error)
}

type linkedListStack struct {
	head  *node
	count int

	lock  *sync.RWMutex
}

func newLinkedListStack() stack {
	stk := &linkedListStack{
		lock: &sync.RWMutex{},
	}

	return stk
}

func (stk *linkedListStack) isEmpty() bool {
	stk.lock.RLock()
	defer stk.lock.RUnlock()

	return stk.count == 0
}

func (stk *linkedListStack) push(index uint64) {
	stk.lock.Lock()
	defer stk.lock.Unlock()

	element := &node{
		data: index,
	}
	temp := stk.head
	element.next = temp
	stk.head = element
	stk.count++
}

func (stk *linkedListStack) pop() (uint64, error) {
	if stk.isEmpty() {
		return 0, errors.New("stack is empty, nothing to pop")
	}

	stk.lock.Lock()
	defer stk.lock.Unlock()

	item := stk.head.data
	stk.head = stk.head.next
	stk.count--

	return item, nil
}
