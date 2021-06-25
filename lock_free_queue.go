package queue

import (
	"sync/atomic"
	"unsafe"
)

type lockFreeQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
	len  int32
}

type node struct {
	value interface{}
	next  unsafe.Pointer
}

func NewLockFreeQueue() *lockFreeQueue {
	n := unsafe.Pointer(&node{})
	return &lockFreeQueue{head: n, tail: n}
}

func (q *lockFreeQueue) Enqueue(v interface{}) {
	n := &node{value: v}
	for {
		tail := load(&q.tail)
		next := load(&tail.next)
		if tail == load(&q.tail) {
			if next == nil {
				if cas(&tail.next, next, n) {
					cas(&q.tail, tail, n)
					atomic.AddInt32(&q.len, 1)
					return
				}
			} else {
				cas(&q.tail, tail, next)
			}
		}
	}
}

func (q *lockFreeQueue) Dequeue() interface{} {
	for {
		head := load(&q.head)
		tail := load(&q.tail)
		next := load(&head.next)
		if head == load(&q.head) {
			if head == tail {
				if next == nil {
					return nil
				}
				cas(&q.tail, tail, next)
			} else {
				v := next.value
				if cas(&q.head, head, next) {
					atomic.AddInt32(&q.len, -1)
					return v
				}
			}
		}
	}
}

func (q *lockFreeQueue) Empty() bool {
	return atomic.LoadInt32(&q.len) != 0
}

func load(p *unsafe.Pointer) *node {
	return (*node)(atomic.LoadPointer(p))
}

func cas(p *unsafe.Pointer, old, new *node) bool {
	return atomic.CompareAndSwapPointer(p,
		unsafe.Pointer(old), unsafe.Pointer(new))
}
