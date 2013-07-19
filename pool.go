package bytepool
// Package bytepool provides a pool of fixed-length []byte

type Pool struct {
  capacity int
  list chan *Item
}

func New(count int, capacity int) *Pool {
  p := &Pool {
    capacity: capacity,
    list: make(chan *Item, count),
  }
  for i := 0; i < count; i++ {
    p.list <- newItem(capacity, p)
  }
  return p
}

func (pool *Pool) Checkout() *Item {
  var item *Item
  select {
    case item = <- pool.list:
    default:
      item = newItem(pool.capacity, nil)
  }
  return item
}

func (pool *Pool) Len() int {
  return len(pool.list)
}
