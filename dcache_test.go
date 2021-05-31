package zcache

import (
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var c1 DCache
var c2 DCache
var c3 DCache
var keys = make([]string, 200)


func randString(n int) string {
	b := make([]byte, 2*n)
	crand.Read(b)
	s := base64.URLEncoding.EncodeToString(b)
	return s[0:n]
}


func TestDCache(t *testing.T) {
	// wait service start
	time.Sleep(2 * time.Second)
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

func init() {
	c1 = NewDCache(
		"localhost:5205",
		[]string{"localhost:5206", "localhost:5207"},
		WithLru(),
		WithSize(1000),
	)
	c2 = NewDCache(
		"localhost:5206",
		[]string{"localhost:5205", "localhost:5207"},
		WithLru(),
		WithSize(1000),
	)
	c3 = NewDCache(
		"localhost:5207",
		[]string{"localhost:5205", "localhost:5206"},
		WithLru(),
		WithSize(1000),
	)
	for i := 0; i < 10000; i++ {
		key := randString(10)
		if i < 200 {
			keys[i] = key
		}
		c1.Set(key, randString(rand.Intn(10000)))
		c2.Set(randString(10), randString(rand.Intn(10000)))
		c3.Set(randString(10), randString(rand.Intn(10000)))
	}
}

func BenchmarkCacheGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c1.Get(keys[i % 200])
	}
}

func BenchmarkCacheSet(b *testing.B) {
	c1.Set(randString(10), randString(rand.Intn(10000)))
}

func TestCost(t *testing.T) {
	v := randString(10000)
	N := 100000
	start := time.Now()
	for i := 0; i < N; i++ {
		c1.Set(randString(10), v)
	}
	fmt.Println("Set QPS", float64(N) / time.Since(start).Seconds())

	start = time.Now()
	for i := 0; i < N; i++ {
		key := randString(10)
		x := rand.Intn(100)
		if x < 10 {
			c1.Set(key, v)
		} else if x < 80 {
			c1.Get(key)
		} else {
			c1.Del(key)
		}
	}
	fmt.Println("Real Use QPS", float64(N) / time.Since(start).Seconds())
}

