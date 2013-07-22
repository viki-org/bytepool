package bytepool

import (
  "io"
)

type Item struct {
  pool *Pool
  length int
  bytes []byte
}

func newItem(capacity int, pool *Pool) *Item {
  return &Item {
    pool: pool,
    bytes: make([]byte, capacity),
  }
}

func (item *Item) Write(b []byte) {
  item.length += copy(item.bytes[item.length:], b)
}

func (item *Item) WriteString(s string) {
  item.length += copy(item.bytes[item.length:], s)
}

func (item *Item) ReadFrom(reader io.Reader) (int64, error) {
  var read int64
  for {
    r, err := reader.Read(item.bytes[item.length:])
    read += int64(r)
    item.length += r
    if err == io.EOF { return read, nil }
    if err != nil { return read, err }
  }
}

func (item *Item) Bytes() []byte {
  return item.bytes[0:item.length]
}

func (item *Item) String() string {
  return string(item.Bytes())
}

func (item *Item) Len() int {
  return item.length;
}

func (item *Item) Position(position int) {
  item.length = position
}

func (item *Item) Close() error{
  item.length = 0
  if item.pool != nil {
    item.pool.list <- item
  }
  return nil
}
