package bytepool

import (
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

func TestCanReadFromAReader(t *testing.T) {
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

func TestCanReadFromMultipleSources(t *testing.T) {
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
