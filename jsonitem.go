package bytepool

import (
  "strconv"
  "strings"
)

type JsonItem struct {
  *Item
  added bool
  pool *JsonPool
  addDelimiter bool
}

func newJsonItem(capacity int, pool *JsonPool) *JsonItem {
  return &JsonItem {
    pool: pool,
    Item: newItem(capacity, nil),
  }
}

var JsonEncode = func (s string) string {
  return strings.Replace(s, `"`, `\"`, -1)
}

func (item *JsonItem) WriteString(s string) int {
  return item.WriteSafeString(JsonEncode(s))
}

func (item *JsonItem) WriteInt(value int) int {
  n := item.Item.WriteString(strconv.Itoa(value))
  return item.delimit(n)
}

func (item *JsonItem) WriteBool(value bool) int {
  n := item.Item.WriteString(strconv.FormatBool(value))
  return item.delimit(n)
}

func (item *JsonItem) WriteSafeString(s string) int {
  return item.writeString(s, true)
}

func (item *JsonItem) WriteKeyString(key, value string) int {
  return item.WriteKeySafeString(key, JsonEncode(value))
}

func (item *JsonItem) WriteKeySafeString(key, value string) int {
  return item.writeKeyValue(key, `"` + value + `"`)
}

func (item *JsonItem) WriteKeyInt(key string, value int) int {
  return item.writeKeyValue(key, strconv.Itoa(value))
}

func (item *JsonItem) WriteKeyBool(key string, value bool) int {
  return item.writeKeyValue(key, strconv.FormatBool(value))
}

func (item *JsonItem) writeKeyValue(key, value string) int {
  n := item.writeString(key, false)
  if item.WriteByte(byte(':')) { n++ }
  n += item.Item.WriteString(value)
  return item.delimit(n)
}

func (item *JsonItem) writeString(s string, delimit bool) int {
  n := item.Item.WriteString(`"` + s + `"`)
  if delimit == false { return n }
  return item.delimit(n)
}

func (item *JsonItem) BeginArray() bool {
  item.added = false
  item.addDelimiter = true
  return item.WriteByte('[')
}

func (item *JsonItem) EndArray() bool {
  item.addDelimiter = false
  if item.added {item.Position(item.Len() - 1)}
  return item.WriteByte(']')
}

func (item *JsonItem) BeginObject() bool {
  item.added = false
  item.addDelimiter = true
  return item.WriteByte('{')
}

func (item *JsonItem) EndObject() bool {
  item.addDelimiter = false
  if item.added {item.Position(item.Len() - 1)}
  return item.WriteByte('}')
}

func (item *JsonItem) delimit(length int) int {
  if item.addDelimiter == false { return length }
  item.WriteByte(byte(','))
  item.added = true
  return length + 1
}

func (item *JsonItem) Close() error {
  item.Item.Close()
  if item.pool != nil {
    item.pool.list <- item
  }
  return nil
}
