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
