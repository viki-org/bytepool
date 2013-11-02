package bytepool

import (
  "io"
  "bytes"
  "testing"
)

func TestCanWriteAString(t *testing.T) {
  expected := "over 9000"
  item := newItem(10, nil)
  item.WriteString("over ")
  item.WriteString("9000")
  actual := item.String()

  if actual != expected {
    t.Errorf("Expecting %q, got %q", expected, actual)
  }
}

func TestCanWriteAByteArray(t *testing.T) {
  expected := []byte("the spice must flow")
  item := newItem(60, nil)
  item.Write([]byte("the "))
  item.Write([]byte("spice "))
  item.Write([]byte("must "))
  item.Write([]byte("flow"))
  actual := item.Bytes()

  if bytes.Compare(actual, expected) != 0 {
    t.Errorf("Expecting %v, got %v", expected, actual)
  }
}

func TestWriteAByte(t *testing.T) {
  expected := []byte("the sp")
  item := newItem(60, nil)
  item.Write([]byte("the "))
  item.WriteByte(byte('s'))
  item.WriteByte(byte('p'))
  actual := item.Bytes()

  if bytes.Compare(actual, expected) != 0 {
    t.Errorf("Expecting %v, got %v", expected, actual)
  }
}

func TestDoesNotWriteAByteWhenFull(t *testing.T) {
  expected := []byte("the s")
  item := newItem(5, nil)
  item.Write([]byte("the "))
  item.WriteByte(byte('s'))
  item.WriteByte(byte('p'))
  actual := item.Bytes()

  if bytes.Compare(actual, expected) != 0 {
    t.Errorf("Expecting %v, got %v", expected, actual)
  }
}

func TestHAndlesReadingAnExactSize(t *testing.T) {
  expected := "12345"
  item := newItem(5, nil)
  buffer := bytes.NewBufferString(expected)
  item.ReadFrom(buffer)

  if item.String() != expected {
    t.Errorf("Expecting %v, got %v", expected, item.String())
  }
}

func TestCanWriteFromAReader(t *testing.T) {
  expected := []byte("I am in a reader")
  item := newItem(60, nil)
  n, _ := item.ReadFrom(bytes.NewBuffer(expected))
  actual := item.Bytes()

  if bytes.Compare(actual, expected) != 0 {
    t.Errorf("Expecting %v, got %v", expected, actual)
  }
  if int(n) != len(expected) {
    t.Errorf("Expecting length of %v, got %v", len(expected), n)
  }
}

func TestCanWriteFromMultipleSources(t *testing.T) {
  expected := []byte("startI am in a readerend")
  bufferContent := []byte("I am in a reader")
  item := newItem(100, nil)
  item.Write([]byte("start"))
  n, _ := item.ReadFrom(bytes.NewBuffer(bufferContent))
  item.WriteString("end")
  actual := item.Bytes()

  if bytes.Compare(actual, expected) != 0 {
    t.Errorf("Expecting %v, got %v", expected, actual)
  }
  if int(n) != len(bufferContent) {
    t.Errorf("Expecting length of %v, got %v", len(bufferContent), n)
  }
}

func TestCanSetThePosition(t *testing.T) {
  expected := "hello."
  item := newItem(100, nil)
  item.WriteString("hello world")
  item.Position(5)
  item.WriteString(".")

  if item.String() != expected {
    t.Errorf("Expecting %v, got %v", expected, item.String())
  }
}

func TestCloseResetsTheLength(t *testing.T) {
  item := newItem(100, nil)
  item.WriteString("hello world")
  item.Close()
  if item.Len() != 0 {
    t.Errorf("Expecting length of 0, got %v", item.Len())
  }
  item.WriteString("hello")

  expected := "hello"
  actual := item.String()
  if actual != expected {
    t.Errorf("Expecting %q, got %q", expected, actual)
  }
}

func TestCannotSetThePositionToANegativeValue(t *testing.T) {
  expected := "hello world."
  item := newItem(25, nil)
  item.WriteString("hello world")
  item.Position(-10)
  item.WriteString(".")

  if item.String() != expected {
    t.Errorf("Expecting %v, got %v", expected, item.String())
  }
}

func TestCannotSetThePositionBeyondTheLength(t *testing.T) {
  expected := "hello world."
  item := newItem(25, nil)
  item.WriteString("hello world")
  item.Position(30)
  item.WriteString(".")

  if item.String() != expected {
    t.Errorf("Expecting %v, got %v", expected, item.String())
  }
}

func TestTrimLastIfTrimsOnMatch(t *testing.T) {
  expected := "hello world"
  item := newItem(25, nil)
  item.WriteString("hello world.")
  r := item.TrimLastIf(byte('.'))
  if r != true {
    t.Error("Expecting TrimLastIf to return true")
  }
  if item.String() != expected {
    t.Errorf("Expecting %v, got %v", expected, item.String())
  }
}

func TestTrimLastIfTrimsOnNoMatch(t *testing.T) {
  expected := "hello world."
  item := newItem(25, nil)
  item.WriteString("hello world.")
  r := item.TrimLastIf(byte(','))
  if r != false {
    t.Error("Expecting TrimLastIf to return false")
  }
  if item.String() != expected {
    t.Errorf("Expecting %v, got %v", expected, item.String())
  }
}

func TestTruncatesTheContentToTheLength(t *testing.T) {
  expected := "hell"
  item := newItem(4, nil)
  item.WriteString("hello")
  actual := item.String()
  if actual != expected {
    t.Errorf("Expecting %q, got %q", expected, actual)
  }

  item.WriteString("world")
  actual = item.String()
  if actual != expected {
    t.Errorf("Expecting %q, got %q", expected, actual)
  }
}

func TestCanReadIntoVariousSizedByteArray(t *testing.T) {
  for size, expected := range map[int]string{3: "hel", 5: "hello", 7: "hello\x00\x00"} {
    item := newItem(5, nil)
    item.WriteString("hello")
    target := make([]byte, size)
    item.Read(target)
    if string(target) != expected {
      t.Errorf("Expecting %q, got %q", expected, string(target))
    }
  }
}

func TestReadDoesNotAutomaticallyRewind(t *testing.T) {
  item := newItem(5, nil)
  item.WriteString("hello")
  b := make([]byte, 5)

  n, err  := item.Read(b[0:2])
  if n != 2 { t.Errorf("expecting to have read 2 bytes, but got %d", n) }
  if err != nil { t.Errorf("should have gotten nil error, got %v", err) }
  if string(b[0:2]) != "he" { t.Errorf("expecting to have read he, got %v", string(b[0:2])) }

  n, err  = item.Read(b[2:])
  if n != 3 { t.Errorf("expecting to have read 3 bytes, but got %d", n)}
  if err != io.EOF { t.Errorf("error should be io.EOF, got %v", err)}
  if string(b[0:5]) != "hello" { t.Errorf("expecting to have read hello, got %v", string(b[0:5])) }

  n, err  = item.Read(b)
  if n != 0 { t.Errorf("expecting to have read 0 bytes, but got %d", n)}
  if err != io.EOF { t.Errorf("error should be io.EOF, got %v", err)}
  if string(b[0:5]) != "hello" { t.Errorf("expecting to have read hello, got %v", string(b[0:5])) }
}
