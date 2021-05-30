package zcache

import (
	"testing"
)

func TestLRU(t *testing.T) {
	cache := newLRUCache(&defaultConfig)
	v := cache.Get("a")
	if v != nil {
		t.Errorf("Get non-nil value from a: %+v", v)
	}
	ev := 123
	cache.Set("a", ev)
	v = cache.Get("a")
	if v != ev {
		t.Errorf("Get value incorrect, expectd=%+v, but get=%+v", ev, v)
	}

	cache.Set("b", 1)
	cache.Set("c", 1)
	cache.Set("d", 1)
	// a invalid?
	v = cache.Get("a")
	if v != nil {
		t.Errorf("Get value incorrect, a is should invalid, but get value=%+v", v)
	}
}



