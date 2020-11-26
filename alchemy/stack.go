package alchemy

import (
	"github.com/go-errors/errors"
	"sync"
)

type node struct {
	data Long
	next *node
}

type stack interface {
	isEmpty() bool
	push(index Long)
	pop() (Long, error)
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

func (stk *linkedListStack) push(index Long) {
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

func (stk *linkedListStack) pop() (Long, error) {
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
