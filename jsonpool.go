package bytepool

import (
	"sync/atomic"
)

// A specialized Pool for making json
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

// This also blocks when there's no more item in pool
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
