package bytepool

// Package bytepool provides a pool of fixed-length []byte

import (
	"sync/atomic"
)

type JsonPool struct {
	misses   int32
	capacity int
	list     chan *JsonItem
}

func NewJson(count int, capacity int) *JsonPool {
	p := &JsonPool{
		capacity: capacity,
		list:     make(chan *JsonItem, count),
	}
	for i := 0; i < count; i++ {
		p.list <- newJsonItem(capacity, p)
	}
	return p
}

func (pool *JsonPool) Checkout() *JsonItem {
	var item *JsonItem
	select {
	case item = <-pool.list:
	default:
		atomic.AddInt32(&pool.misses, 1)
		item = newJsonItem(pool.capacity, nil)
	}
	return item
}

func (pool *JsonPool) Len() int {
	return len(pool.list)
}

func (pool *JsonPool) Misses() int32 {
	return atomic.LoadInt32(&pool.misses)
}
