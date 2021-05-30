package zcache

import (
	"sync"
)

type List struct {
	head *Item
	len int
}

func newList() *List {
	l := &List{
		head: &Item{},
		len: 0,
	}
	l.head.next = l.head
	l.head.prev = l.head
	return l
}

func (l *List) remove(i *Item) *Item {
	i.prev.next = i.next
	if i.next != nil {
		i.next.prev = i.prev
	}
	// i is tail?
	if l.isTail(i) {
		l.head.prev = i.prev
	}
	i.next = nil
	i.prev = nil
	l.len--
	return i
}

// return old tail
func (l *List) transferTail() *Item {
	if l.isTail(l.head) {
		return nil
	}
	return l.transfer(l.head.prev)
}


func (l *List) isTail(i *Item) bool {
	return l.head.prev == i
}

func (l *List) insert(i *Item) *Item {
	if l.head.next != nil {
		l.head.next.prev = i
	}
	l.head.next = i
	if l.isTail(l.head) {
		l.head.prev = i
	}
	l.len++
	return i
}

// insert to head
func (l *List) insertValue(k string, v interface{}) *Item {
	i := &Item{
		next: l.head.next,
		prev: l.head,
		value: v,
		key: k,
	}
	return l.insert(i)
}

// i be removed and then insert to head
func (l *List) transfer(i *Item) *Item {
	if i.prev == l.head {
		return i
	}
	if l.isTail(i) {
		l.head.prev = i.prev
	}
	// removed
	i.prev.next = i.next
	if i.next != nil {
		i.next.prev = i.prev
	}
	// insert
	i.next = l.head.next
	i.prev = l.head
	l.head.next = i
	return i
}


type Item struct {
	key   string
	value interface{}		// set(k, v) v will save in here
	prev *Item
	next *Item
}

func (i *Item) isInvalid() bool {
	return i.prev == nil || i.next == nil		// and expired?
}

type LRUCache struct {
	mu sync.RWMutex
	list   	  *List
	items     map[string]*Item
	keySize   int64			// ignore
	valueSize int64			// ignore
	limitSize int64			// list len
}

func newLRUCache() *LRUCache {
	return &LRUCache{
		list: newList(),
		items: map[string]*Item{},
		limitSize: 3,
	}
}


func (c *LRUCache) Get(key string) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.items[key]
	if !ok {
		return nil
	}
	c.list.transfer(item)
	return item.value
}

func (c *LRUCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, ok := c.items[key]; ok {
		item.value = value
		c.list.transfer(item)
	} else {
		if int64(c.list.len) == c.limitSize {
			item = c.list.transferTail()
			delete(c.items, item.key)
			item.key = key
			item.value = value
		} else {
			item = c.list.insertValue(key, value)
		}
		c.items[key] = item
	}
}

func (c *LRUCache) Del(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, ok := c.items[key]; ok {
		c.list.remove(item)
		delete(c.items, key)
	}
}


func (c *LRUCache) clearInvalid() {
	if len(c.items) <= c.list.len {
		return
	}

}


