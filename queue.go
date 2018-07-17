package fastqueue

import (
	"runtime"
	"sync/atomic"
)

const CacheLinePaddingSize = 7

type element struct {
	pos  uint64
	data interface{}
}

type Queue struct {
	head     uint64
	p1       [CacheLinePaddingSize]uint64
	tail     uint64
	p2       [CacheLinePaddingSize]uint64
	mask     uint64
	p3       [CacheLinePaddingSize]uint64
	elements []*element
}

func NewQueue(size uint64) *Queue {
	size = roundUp(size) // round up the size to next power of 2

	q := &Queue{}
	q.elements = make([]*element, size, size)
	for i := uint64(0); i < size; i++ {
		q.elements[i] = &element{pos: i}
	}
	q.mask = size - 1

	return q
}

func (q *Queue) Push(item interface{}) bool {
	var head uint64

	for {
		head = atomic.LoadUint64(&q.head)

		if atomic.CompareAndSwapUint64(&q.head, head, head+1) {
			break
		}

		runtime.Gosched()
	}

	e := q.elements[head&q.mask]
	e.data = item

	return true
}

func (q *Queue) Pop() interface{} {
	var tail uint64

	for {
		tail = atomic.LoadUint64(&q.tail)

		if atomic.CompareAndSwapUint64(&q.tail, tail, tail+1) {
			break
		}

		runtime.Gosched()
	}

	e := q.elements[tail&q.mask]
	return e.data
}
