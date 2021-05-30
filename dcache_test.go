package zcache

import (
	"fmt"
	"testing"
	"time"
)

func TestDCache(t *testing.T) {
	c1:= NewDCache(
		"localhost:5205",
		[]string{"localhost:5206", "localhost:5207"},
		WithLru(),
	)
	c2:= NewDCache(
		"localhost:5206",
		[]string{"localhost:5205", "localhost:5207"},
		WithLru(),
	)
	c3:= NewDCache(
		"localhost:5207",
		[]string{"localhost:5205", "localhost:5206"},
		WithLru(),
	)
	// wait service start
	time.Sleep(2 * time.Second)
	fmt.Println(c1, c2)
	c1.Set("a", 1)
	c2.Set("a", 2)
	c3.Set("a", 2)
	fmt.Println("c2 Get(a)", c2.Get("a"))
	c1.Del("a")
	c1.Del("d")
	fmt.Println("c2 Get(a) after c1 del", c2.Get("a"))
	fmt.Println("c3 Get(a) after c1 del", c3.Get("a"))

	time.Sleep(105 * time.Millisecond)
	fmt.Println("c2 Get(a) after sleep", c2.Get("a"))
	fmt.Println("c3 Get(a) after sleep", c3.Get("a"))
}



