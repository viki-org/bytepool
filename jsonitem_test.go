package bytepool

import (
  "testing"
)

func TestJsonCanWriteAnEncodedString(t *testing.T) {
  expected := `"over \"9000\""`
  item := newJsonItem(100, nil)
  item.WriteString(`over "9000"`)
  actual := item.String()
  if actual != expected {
    t.Errorf("Expecting %q, got %q", expected, actual)
  }
}

func TestJsonCanWriteAString(t *testing.T) {
  expected := `"over "9000""`
  item := newJsonItem(100, nil)
  item.WriteSafeString(`over "9000"`)
  actual := item.String()
  if actual != expected {
    t.Errorf("Expecting %q, got %q", expected, actual)
  }
}

func TestJsonWritesAnEmptyArray(t *testing.T) {
  expected := "[]"
  item := newJsonItem(100, nil)
  item.BeginArray()
  item.EndArray()
  actual := item.String()
  if actual != expected {
    t.Errorf("Expecting %q, got %q", expected, actual)
  }
}

func TestJsonWritesASingleValueArray(t *testing.T) {
  expected := "[90]"
  item := newJsonItem(100, nil)
  item.BeginArray()
  item.WriteInt(90)
  item.EndArray()
  actual := item.String()
  if actual != expected {
    t.Errorf("Expecting %q, got %q", expected, actual)
  }
}

func TestJsonWritesAMultiValueArray(t *testing.T) {
  expected := `[90,false,"abc",true]`
  item := newJsonItem(100, nil)
  item.BeginArray()
  item.WriteInt(90)
  item.WriteBool(false)
  item.WriteString("abc")
  item.WriteBool(true)
  item.EndArray()
  actual := item.String()
  if actual != expected {
    t.Errorf("Expecting %q, got %q", expected, actual)
  }
}

func TestJsonWritesAnEmptyObject(t *testing.T) {
  expected := `{}`
  item := newJsonItem(100, nil)
  item.BeginObject()
  item.EndObject()
  actual := item.String()
  if actual != expected {
    t.Errorf("Expecting %q, got %q", expected, actual)
  }
}

func TestJsonASingleValueObject(t *testing.T) {
  expected := `{"over":"90\"00!"}`
  item := newJsonItem(100, nil)
  item.BeginObject()
  item.WriteKeyString("over", "90\"00!")
  item.EndObject()
  actual := item.String()
  if actual != expected {
    t.Errorf("Expecting %q, got %q", expected, actual)
  }
}

func TestJsonAMultiValueObject(t *testing.T) {
  expected := `{"name":"goku","power":9000,"over":true}`
  item := newJsonItem(100, nil)
  item.BeginObject()
  item.WriteKeySafeString("name", "goku")
  item.WriteKeyInt("power", 9000)
  item.WriteKeyBool("over", true)
  item.EndObject()
  actual := item.String()
  if actual != expected {
    t.Errorf("Expecting %q, got %q", expected, actual)
  }
}
