package bytepool

import (
  "bytes"
  "io"
  . "gopkg.in/check.v1"
)

func (s *TestSuite) TestCanWriteAString(c *C) {
  expected := "over 9000"
  item := newItem(10, nil)
  item.WriteString("over ")
  item.WriteString("9000")
  actual := item.String()

  c.Assert(actual, Equals, expected, Commentf("Expecting %q, got %q", expected, actual))
}

func (s *TestSuite) TestCanWriteAByteArray(c *C) {
  expected := []byte("the spice must flow")
  item := newItem(60, nil)
  item.Write([]byte("the "))
  item.Write([]byte("spice "))
  item.Write([]byte("must "))
  item.Write([]byte("flow"))
  actual := item.Bytes()

  c.Assert(actual, DeepEquals, expected, Commentf("Expecting %q, got %q", expected, actual))
}

func (s *TestSuite) TestWriteAByte(c *C) {
  expected := []byte("the sp")
  item := newItem(60, nil)
  item.Write([]byte("the "))
  item.WriteByte(byte('s'))
  item.WriteByte(byte('p'))
  actual := item.Bytes()

  c.Assert(actual, DeepEquals, expected, Commentf("Expecting %q, got %q", expected, actual))
}

func (s *TestSuite) TestDoesNotWriteAByteWhenFull(c *C) {
  expected := []byte("the s")
  item := newItem(5, nil)
  item.Write([]byte("the "))
  item.WriteByte(byte('s'))
  item.WriteByte(byte('p'))
  actual := item.Bytes()

  c.Assert(actual, DeepEquals, expected, Commentf("Expecting %q, got %q", expected, actual))
}

func (s *TestSuite) TestHAndlesReadingAnExactSize(c *C) {
  expected := "12345"
  item := newItem(5, nil)
  buffer := bytes.NewBufferString(expected)
  item.ReadFrom(buffer)

  c.Assert(item.String(), Equals, expected, Commentf("Expecting %v, got %v", expected, item.String()))
}

func (s *TestSuite) TestCanWriteFromAReader(c *C) {
  expected := []byte("I am in a reader")
  item := newItem(60, nil)
  n, _ := item.ReadFrom(bytes.NewBuffer(expected))
  actual := item.Bytes()

  c.Assert(actual, DeepEquals, expected, Commentf("Expecting %q, got %q", expected, actual))
  c.Assert(int(n), Equals, len(expected), Commentf("Expecting %v, got %v", len(expected), int(n)))
}

func (s *TestSuite) TestCanWriteFromMultipleSources(c *C) {
  expected := []byte("startI am in a readerend")
  bufferContent := []byte("I am in a reader")
  item := newItem(100, nil)
  item.Write([]byte("start"))
  n, _ := item.ReadFrom(bytes.NewBuffer(bufferContent))
  item.WriteString("end")
  actual := item.Bytes()

  c.Assert(actual, DeepEquals, expected, Commentf("Expecting %q, got %q", expected, actual))
  c.Assert(int(n), Equals, len(bufferContent), Commentf("Expecting %v, got %v", len(bufferContent), int(n)))
}

func (s *TestSuite) TestCanSetThePosition(c *C) {
  expected := "hello."
  item := newItem(100, nil)
  item.WriteString("hello world")
  item.Position(5)
  item.WriteString(".")

  c.Assert(item.String(), Equals, expected, Commentf("Expecting %v, got %v", expected, item.String()))
}

func (s *TestSuite) TestCloseResetsTheLength(c *C) {
  item := newItem(100, nil)
  item.WriteString("hello world")
  item.Close()

  c.Assert(item.Len(), Equals, 0, Commentf("Expecting length of 0, got %v", item.Len()))

  item.WriteString("hello")
  expected := "hello"
  actual := item.String()
  c.Assert(actual, Equals, expected, Commentf("Expecting %q, got %q", expected, actual))
}

func (s *TestSuite) TestCannotSetThePositionToANegativeValue(c *C) {
  expected := "hello world."
  item := newItem(25, nil)
  item.WriteString("hello world")
  item.Position(-10)
  item.WriteString(".")

  c.Assert(item.String(), Equals, expected, Commentf("Expecting %v, got %v", expected, item.String()))
}

func (s *TestSuite) TestCannotSetThePositionBeyondTheLength(c *C) {
  expected := "hello world."
  item := newItem(25, nil)
  item.WriteString("hello world")
  item.Position(30)
  item.WriteString(".")

  c.Assert(item.String(), Equals, expected, Commentf("Expecting %v, got %v", expected, item.String()))
}

func (s *TestSuite) TestTrimLastIfTrimsOnMatch(c *C) {
  expected := "hello world"
  item := newItem(25, nil)
  item.WriteString("hello world.")
  r := item.TrimLastIf(byte('.'))
  c.Assert(r, Equals, true, Commentf("Expecting TrimLastIf to return true"))
  c.Assert(item.String(), Equals, expected, Commentf("Expecting %v, got %v", expected, item.String()))
}

func (s *TestSuite) TestTrimLastIfTrimsOnNoMatch(c *C) {
  expected := "hello world."
  item := newItem(25, nil)
  item.WriteString("hello world.")
  r := item.TrimLastIf(byte(','))
  c.Assert(r, Equals, false, Commentf("Expecting TrimLastIf to return true"))
  c.Assert(item.String(), Equals, expected, Commentf("Expecting %v, got %v", expected, item.String()))
}

func (s *TestSuite) TestTruncatesTheContentToTheLength(c *C) {
  expected := "hell"
  item := newItem(4, nil)
  item.WriteString("hello")
  actual := item.String()

  c.Assert(actual, Equals, expected, Commentf("Expecting %q, got %q", expected, actual))

  item.WriteString("world")
  actual = item.String()

  c.Assert(actual, Equals, expected, Commentf("Expecting %q, got %q", expected, actual))
}

func (s *TestSuite) TestCanReadIntoVariousSizedByteArray(c *C) {
  for size, expected := range map[int]string{3: "hel", 5: "hello", 7: "hello\x00\x00"} {
    item := newItem(5, nil)
    item.WriteString("hello")
    target := make([]byte, size)
    item.Read(target)
    c.Assert(string(target), Equals, expected, Commentf("Expecting %q, got %q", expected, string(target)))
  }
}

func (s *TestSuite) TestReadDoesNotAutomaticallyRewind(c *C) {
  item := newItem(5, nil)
  item.WriteString("hello")
  b := make([]byte, 5)

  n, err := item.Read(b[0:2])

  c.Assert(n, Equals, 2, Commentf("expecting to have read 2 bytes, but got %d", n))
  c.Assert(err, IsNil, Commentf("should have gotten nil error, got %v", err))
  c.Assert(string(b[0:2]), Equals, "he", Commentf("expecting to have read `he`, got %v", string(b[0:2])))

  n, err = item.Read(b[2:])

  c.Assert(n, Equals, 3, Commentf("expecting to have read 3 bytes, but got %d", n))
  c.Assert(err, Equals, io.EOF, Commentf("error should be io.EOF, got %v", err))
  c.Assert(string(b[0:5]), Equals, "hello", Commentf("expecting to have read `hello`, got %v", string(b[0:5])))

  n, err = item.Read(b)

  c.Assert(n, Equals, 0, Commentf("expecting to have read 0 bytes, but got %d", n))
  c.Assert(err, Equals, io.EOF, Commentf("error should be io.EOF, got %v", err))
  c.Assert(string(b[0:5]), Equals, "hello", Commentf("expecting to have read `hello`, got %v", string(b[0:5])))
}
