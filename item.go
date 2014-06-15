package bytepool

import (
	"io"
)

// a slice of bytes within the pool
//    pool: points to the pool containing this item
//    length:
//    read:
//    bytes: the slice
type Item struct {
	pool   *Pool
	length int
	read   int
	bytes  []byte
}

func newItem(capacity int, pool *Pool) *Item {
	return &Item{
		pool:  pool,
		bytes: make([]byte, capacity),
	}
}

func (item *Item) Write(b []byte) (int, error) {
	if item.Full() {
		return 0, io.ErrShortWrite
	}
	n := copy(item.bytes[item.length:], b)
	item.length += n
	return n, nil
}

func (item *Item) WriteByte(b byte) bool {
	if item.Full() {
		return false
	}
	item.bytes[item.length] = b
	item.length += 1
	return true
}

func (item *Item) WriteString(s string) int {
	if item.Full() {
		return 0
	}
	n := copy(item.bytes[item.length:], s)
	item.length += n
	return n
}

// read data from an io.Reader into the item's slice
func (item *Item) ReadFrom(reader io.Reader) (int64, error) {
	var read int64
	for {
		r, err := reader.Read(item.bytes[item.length:])
		read += int64(r)
		item.length += r
		if err == io.EOF || item.Full() {
			return read, nil
		}
		if err != nil {
			return read, err
		}
	}
}

// read data from the item in to another byte-slice
func (item *Item) Read(p []byte) (int, error) {
	if item.Drained() {
		return 0, io.EOF
	}
	n := copy(p, item.bytes[item.read:item.length])
	item.read += n
	if item.Drained() {
		return n, io.EOF
	}
	return n, nil
}

// return only the content that has been read so far
func (item *Item) Bytes() []byte {
	return item.bytes[0:item.length]
}

// return the full slice
func (item *Item) Raw() []byte {
	return item.bytes
}

func (item *Item) String() string {
	return string(item.Bytes())
}

func (item *Item) Len() int {
	return item.length
}

// trim last item of a non-empty slice that's equal to the param
func (item *Item) TrimLastIf(b byte) bool {
	l := item.length - 1
	if l == -1 || item.bytes[l] != b {
		return false
	}
	item.Position(l)
	return true
}

// update the size of the safe-to-read subset of the slice
func (item *Item) Position(position int) bool {
	if position < 0 || position > cap(item.bytes) {
		return false
	}
	item.length = position
	return true
}

func (item *Item) Full() bool {
	return item.length == cap(item.bytes)
}

// tell whether all the content in the slice has been read
func (item *Item) Drained() bool {
	return item.length == item.read
}

// close the item and return it to the pool
func (item *Item) Close() error {
	item.length = 0
	item.read = 0
	if item.pool != nil {
		item.pool.list <- item
	}
	return nil
}
