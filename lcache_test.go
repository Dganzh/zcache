package zcache

import (
	"fmt"
	"testing"
)

func TestLCache(t *testing.T) {
	cache := NewCache(
		WithLru(),
	)
	fmt.Println(cache)
	cache.Get("a")
}



