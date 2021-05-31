package zcache

import (
	"testing"
)

func TestLCache(t *testing.T) {
	cache := NewCache()
	v := 123
	cache.Set("a", v)
	if v != cache.Get("a") {
		t.Errorf("Get invalid value: %+v", v)
	}
	cache.Del("a")
	if cache.Get("a") != nil {
		t.Error("Del failed")
	}
}
