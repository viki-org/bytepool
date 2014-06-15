package bytepool

import (
	"strconv"
	"strings"
	"time"
)

// an item specially designed for writing JSON, inherrits Item's behavior
//    depth: indicate nesting level (object & array in json), start at 0
//    added: might be an unnecessary field
//    pool: the pool of json items
type JsonItem struct {
	*Item
	depth int
	added bool
	pool  *JsonPool
}

func newJsonItem(capacity int, pool *JsonPool) *JsonItem {
	return &JsonItem{
		pool: pool,
		Item: newItem(capacity, nil),
	}
}

var JsonEncode = strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace

// Write arbitary string into the slice
func (item *JsonItem) WriteString(s string) int {
	return item.WriteSafeString(JsonEncode(s))
}

// Write arbitary int into the slice
func (item *JsonItem) WriteInt(value int) int {
	n := item.Item.WriteString(strconv.Itoa(value))
	return item.delimit(n)
}

// Write bool value into the slice
func (item *JsonItem) WriteBool(value bool) int {
	n := item.Item.WriteString(strconv.FormatBool(value))
	return item.delimit(n)
}

// Write time value (in RFC3339, or yyyy-mm-ddThh:mm:ssZ
func (item *JsonItem) WriteTime(value time.Time) int {
	return item.WriteString(value.Format(time.RFC3339))
}

// Write an safe (doesn't need escaping) string into the slice
func (item *JsonItem) WriteSafeString(s string) int {
	return item.writeString(s, true)
}



// Write a key-value pair where value is non-empty string
func (item *JsonItem) WriteKeySafeString(key, value string) int {
	return item.WriteKeyValue(key, `"`+value+`"`)
}

// Write a key-value pair where value is any string (might need escaping)
func (item *JsonItem) WriteKeyString(key, value string) int {
	return item.WriteKeySafeString(key, JsonEncode(value))
}

// Write a key-value pair where value is an int
func (item *JsonItem) WriteKeyInt(key string, value int) int {
	return item.WriteKeyValue(key, strconv.Itoa(value))
}

// Write a key-value pair where value is bool
func (item *JsonItem) WriteKeyBool(key string, value bool) int {
	return item.WriteKeyValue(key, strconv.FormatBool(value))
}

// Write a key-value pair where value is a time (in yyyy-mm-ddThh-mm-ssZ)
func (item *JsonItem) WriteKeyTime(key string, value time.Time) int {
	return item.WriteKeyString(key, value.Format(time.RFC3339))
}

// Write a partial(to `[` char) key-value pair where value is an array
func (item *JsonItem) WriteKeyArray(key string) int {
	n := item.writeString(key, false)
	if item.WriteByte(byte(':')) {
		n++
	}
	if item.BeginArray() {
		n++
	}
	return n
}

// Write a partial(to `{` char) key-value pair where value is an object
func (item *JsonItem) WriteKeyObject(key string) int {
	n := item.writeString(key, false)
	if item.WriteByte(byte(':')) {
		n++
	}
	if item.BeginObject() {
		n++
	}
	return n
}

// Write a key-value pair where value is any string-casted values
func (item *JsonItem) WriteKeyValue(key, value string) int {
	n := item.writeString(key, false)
	if item.WriteByte(byte(':')) {
		n++
	}
	n += item.Item.WriteString(value)
	return item.delimit(n)
}


// Write a string & quotes into the slice.Typically used to write keys
func (item *JsonItem) writeString(s string, delimit bool) int {
	n := item.Item.WriteString(`"` + s + `"`)
	if delimit == false {
		return n
	}
	return item.delimit(n)
}

// Write Array start: `[` into the slice and increase depth
func (item *JsonItem) BeginArray() bool {
	item.added = false
	item.depth++
	return item.WriteByte('[')
}

// Write Array end: `]` into the slice and decrease depth
// also trim extra comma at the end
func (item *JsonItem) EndArray() (int, error) {
	item.depth--
	item.TrimLastIf(',')
	return item.Write([]byte("],"))
}

// Write Object start: `{` into the slice and increase depth
func (item *JsonItem) BeginObject() bool {
	item.depth++
	return item.WriteByte('{')
}

// Write Object start: `}` into the slice and decrease depth
// also trim extra comma at the end
func (item *JsonItem) EndObject() (int, error) {
	item.depth--
	item.TrimLastIf(',')
	return item.Write([]byte("},"))
}

// Write an array / object delimiter: `,` into the slice
func (item *JsonItem) delimit(length int) int {
	if item.depth == 0 {
		return length
	}
	item.WriteByte(',')
	return length + 1
}

// Close the JsonItem and return it to the pool
func (item *JsonItem) Close() error {
	item.TrimLastIf(',')
	item.Item.Close()
	if item.pool != nil {
		item.depth = 0
		item.pool.list <- item
	}
	return nil
}

func (item *JsonItem) Len() int {
	item.TrimLastIf(',')
	return item.Item.Len()
}

// Return the meaningful (intentionally created) bytes in the JsonItem
func (item *JsonItem) Bytes() []byte {
	item.TrimLastIf(',')
	return item.Item.Bytes()
}

// Return the whole slice (including jumble bytes)
func (item *JsonItem) Raw() []byte {
	item.TrimLastIf(',')
	return item.Item.Raw()
}

// Return the string made from the slice
func (item *JsonItem) String() string {
	item.TrimLastIf(',')
	return item.Item.String()
}
