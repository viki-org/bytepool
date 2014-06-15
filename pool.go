/*
Create a pool of fixed-length byte-slices shared between go routines
By caching unused slices, it reduce garbage-collector work
hence faster programs

This is similar to sync.Pool in (unreleased) Go 1.3
*/
package bytepool

// Package bytepool provides a pool of fixed-length []byte

import (
	"sync/atomic"
)

// The pool of byte-slices
//    misses: count when checkout fails (there's no more slices)
//    capacity: size of each slices 
//    list: the pool
type Pool struct {
	misses   int32
	capacity int
	list     chan *Item
}

func New(count int, capacity int) *Pool {
	p := &Pool{
		capacity: capacity,
		list:     make(chan *Item, count),
	}
	for i := 0; i < count; i++ {
		p.list <- newItem(capacity, p)
	}
	return p
}

// Get an item out from the pool
// when there are not enough slices available, it blocks
// and also increase the misses count
func (pool *Pool) Checkout() *Item {
	var item *Item
	select {
	case item = <-pool.list:
	default:
		atomic.AddInt32(&pool.misses, 1)
		item = newItem(pool.capacity, nil)
	}
	return item
}

// no of items left inside the pool
func (pool *Pool) Len() int {
	return len(pool.list)
}

func (pool *Pool) Misses() int32 {
	return atomic.LoadInt32(&pool.misses)
}
