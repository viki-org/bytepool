package bytepool

import (
  "time"
  . "gopkg.in/check.v1"
)

func (s *TestSuite) TestJsonCanWriteAnEncodedString(c *C) {
  expected := `"over \"9000\""`
  item := newJsonItem(100, nil)
  item.WriteString(`over "9000"`)
  actual := item.String()

  c.Assert(actual, Equals, expected, Commentf("Expecting %q, got %q", expected, actual))
}

func (s *TestSuite) TestJsonCanWriteAnEncodedStringWithSlashes(c *C) {
  expected := `"\\over \"\\\\9000/\""`
  item := newJsonItem(100, nil)
  item.WriteString(`\over "\\9000/"`)
  actual := item.String()
  c.Assert(actual, Equals, expected, Commentf("Expecting %q, got %q", expected, actual))

}

func (s *TestSuite) TestJsonCanWriteAString(c *C) {
  expected := `"over "9000""`
  item := newJsonItem(100, nil)
  item.WriteSafeString(`over "9000"`)
  actual := item.String()

  c.Assert(actual, Equals, expected, Commentf("Expecting %q, got %q", expected, actual))
}

func (s *TestSuite) TestJsonWritesAnEmptyArray(c *C) {
  expected := "[]"
  item := newJsonItem(100, nil)
  item.BeginArray()
  item.EndArray()
  actual := item.String()

  c.Assert(actual, Equals, expected, Commentf("Expecting %q, got %q", expected, actual))
}

func (s *TestSuite) TestJsonWritesASingleValueArray(c *C) {
  expected := "[90]"
  item := newJsonItem(100, nil)
  item.BeginArray()
  item.WriteInt(90)
  item.EndArray()
  actual := item.String()

  c.Assert(actual, Equals, expected, Commentf("Expecting %q, got %q", expected, actual))
}

func (s *TestSuite) TestJsonWritesAMultiValueArray(c *C) {
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

  c.Assert(actual, Equals, expected, Commentf("Expecting %q, got %q", expected, actual))
}

func (s *TestSuite) TestJsonWritesAnEmptyObject(c *C) {
  expected := `{}`
  item := newJsonItem(100, nil)
  item.BeginObject()
  item.EndObject()
  actual := item.String()

  c.Assert(actual, Equals, expected, Commentf("Expecting %q, got %q", expected, actual))
}

func (s *TestSuite) TestJsonASingleValueObject(c *C) {
  expected := `{"over":"90\"00!"}`
  item := newJsonItem(100, nil)
  item.BeginObject()
  item.WriteKeyString("over", "90\"00!")
  item.EndObject()
  actual := item.String()

  c.Assert(actual, Equals, expected, Commentf("Expecting %q, got %q", expected, actual))
}

func (s *TestSuite) TestJsonAMultiValueObject(c *C) {
  expected := `{"name":"goku","power":9000,"over":true,"time":"2012-12-12T00:00:00Z"}`
  item := newJsonItem(100, nil)
  item.BeginObject()
  item.WriteKeySafeString("name", "goku")
  item.WriteKeyInt("power", 9000)
  item.WriteKeyBool("over", true)
  item.WriteKeyTime("time", time.Date(2012, time.December, 12, 0, 0, 0, 0, time.UTC))
  item.EndObject()
  actual := item.String()

  c.Assert(actual, Equals, expected, Commentf("Expecting %q, got %q", expected, actual))
}

func (s *TestSuite) TestJsonNestedObjects(c *C) {
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

  c.Assert(actual, Equals, expected, Commentf("Expecting %q, got %q", expected, actual))
}
