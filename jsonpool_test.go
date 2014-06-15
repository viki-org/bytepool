package bytepool

import (
	"reflect"
  . "gopkg.in/check.v1"
)

func (s *TestSuite) TestJsonPoolEachItemIsOfASpecifiedSize(c *C) {
	p := New(1, 9)
	item := p.Checkout()
	defer item.Close()

  c.Assert(cap(item.bytes), Equals, 9, Commentf("expecting array to have a capacity of %d, got %d", 9, cap(item.bytes)))
}

func (s *TestSuite) TestJsonPoolDynamicallyCreatesAnItemWhenPoolIsEmpty(c *C) {
	p := New(1, 2)
	item1 := p.Checkout()
	item2 := p.Checkout()

  c.Assert(cap(item2.bytes), Equals, 2, Commentf("Dynamically created item was not properly initialized"))
  c.Assert(item2.pool, IsNil, Commentf("The dynamically created item should have a nil pool"))

	item1.Close()
	item2.Close()

  c.Assert(p.Len(), Equals, 1, Commentf("Expecting a pool length of 1, got %d", p.Len()))
  c.Assert(p.Misses(), Equals, 1, Commentf("Expecting a miss count of 1, got %d", p.Misses()))
}

func (s *TestSuite) TestJsonPoolReleasesAnItemBackIntoThePool(c *C) {
	p := New(1, 20)
	item1 := p.Checkout()
	pointer := reflect.ValueOf(item1).Pointer()
	item1.Close()
	item2 := p.Checkout()
	defer item2.Close()

  c.Assert(reflect.ValueOf(item2).Pointer(), Equals, pointer, Commentf("Pool returned an unexpected item"))
}
