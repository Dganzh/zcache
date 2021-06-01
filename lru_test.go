package zcache

import (
	"fmt"
	"sync"
	"testing"
	"time"
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

func TestLoad(t *testing.T) {
	cache := newLRUCache(&defaultConfig)
	queryDb := func() (interface{}, error) {
		time.Sleep(50 * time.Millisecond)
		fmt.Println("query in db")
		return "abc", nil	//errors.New("db error")
	}

	var wg sync.WaitGroup
	get := func() {
		v, e := cache.Load("a", queryDb)
		fmt.Println("get result v", v, "e", e)
		wg.Done()
	}

	n := 10
	wg.Add(n)
	// 前一半直接并发取missing key数据
	// 后半部分等数据库查询结束后再发起
	for i := 0; i < n; i++ {
		go get()
		if i == n / 2 {
			time.Sleep(55 * time.Millisecond)
		}
	}
	wg.Wait()
}

func TestExpire(t *testing.T) {
	exp := 5 * time.Second
	cfg := &Config{
		evictType: EvictLru,
		size: 3,
		expire: &exp,
	}
	cache := newLRUCache(cfg)
	cache.SetWithExpire("a", 123, 100 * time.Millisecond)
	cache.SetWithExpire("b", 123, 100 * time.Millisecond)
	cache.Set("c", 123)
	fmt.Println(cache.Get("a") == 123)
	time.Sleep(100 * time.Millisecond)
	// Get 一个过期key会触发删除
	fmt.Println(cache.Get("a") == nil)

	// 这里只剩下b c
	keys := cache.Keys()
	fmt.Println("keys", keys, "len", len(keys))
	// 确保触发周期删除过期key
	time.Sleep(1 * time.Second)
	keys = cache.Keys()
	fmt.Println("after auto delete, keys:", keys, "len", len(keys))

	fmt.Println(cache.Get("c") == 123)
	time.Sleep(exp)
	fmt.Println(cache.Get("c") == nil)
}
