package bytepool

import (
	"testing"
	"time"
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

func TestJsonCanWriteAnEncodedStringWithSlashes(t *testing.T) {
	expected := `"\\over \"\\\\9000/\""`
	item := newJsonItem(100, nil)
	item.WriteString(`\over "\\9000/"`)
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
	expected := `[90,false,"2012-12-12T00:00:00Z","abc",true]`
	item := newJsonItem(100, nil)
	item.BeginArray()
	item.WriteInt(90)
	item.WriteBool(false)
	item.WriteTime(time.Date(2012, time.December, 12, 0, 0, 0, 0, time.UTC))
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
	expected := `{"name":"goku","power":9000,"over":true,"time":"2012-12-12T00:00:00Z"}`
	item := newJsonItem(100, nil)
	item.BeginObject()
	item.WriteKeySafeString("name", "goku")
	item.WriteKeyInt("power", 9000)
	item.WriteKeyBool("over", true)
	item.WriteKeyTime("time", time.Date(2012, time.December, 12, 0, 0, 0, 0, time.UTC))
	item.EndObject()
	actual := item.String()
	if actual != expected {
		t.Errorf("Expecting %q, got %q", expected, actual)
	}
}

func TestJsonNestedObjects(t *testing.T) {
	expected := `[1,{"name":"goku","levels":[2,{"over":{"9000":"!"},"acquired":"2012-12-12T00:00:00Z"}],"age":12}]`
	item := newJsonItem(100, nil)
	item.BeginArray()
	item.WriteInt(1)
	item.BeginObject()
	item.WriteKeyString("name", "goku")
	item.WriteKeyArray("levels")
	item.WriteInt(2)
	item.BeginObject()
	item.WriteKeyObject("over")
	item.WriteKeyString("9000", "!")
	item.EndObject()
	item.WriteKeyTime("acquired", time.Date(2012, time.December, 12, 0, 0, 0, 0, time.UTC))
	item.EndObject()
	item.EndArray()
	item.WriteKeyInt("age", 12)
	item.EndObject()
	item.EndArray()
	actual := item.String()
	if actual != expected {
		t.Errorf("Expecting %q, got %q", expected, actual)
	}
}
