package bytepool

import (
	"reflect"
	"testing"
)

func TestPoolEachItemIsOfASpecifiedSize(t *testing.T) {
	expected := 9
	p := New(1, expected)
	item := p.Checkout()
	defer item.Close()
	if cap(item.bytes) != expected {
		t.Errorf("expecting array to have a capacity of %d, got %d", expected, cap(item.bytes))
	}
}

func TestPoolDynamicallyCreatesAnItemWhenPoolIsEmpty(t *testing.T) {
	p := New(1, 2)
	item1 := p.Checkout()
	item2 := p.Checkout()
	if cap(item2.bytes) != 2 {
		t.Error("Dynamically created item was not properly initialized")
	}
	if item2.pool != nil {
		t.Error("The dynamically created item should have a nil pool")
	}

	item1.Close()
	item2.Close()
	if p.Len() != 1 {
		t.Errorf("Expecting a pool lenght of 1, got %d", p.Len())
	}
	if p.Misses() != 1 {
		t.Errorf("Expecting a miss count of 1, got %d", p.Misses())
	}

}
func TestPoolReleasesAnItemBackIntoThePool(t *testing.T) {
	p := New(1, 20)
	item1 := p.Checkout()
	pointer := reflect.ValueOf(item1).Pointer()
	item1.Close()

	item2 := p.Checkout()
	defer item2.Close()
	if reflect.ValueOf(item2).Pointer() != pointer {
		t.Error("Pool returned an unexected item")
	}
}
