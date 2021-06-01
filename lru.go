package zcache

import (
	"sync"
	"time"
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
	expire *time.Time
	prev *Item
	next *Item
}

func (i *Item) isExpired() bool {
	if i.expire == nil {
		return false
	}
	return i.expire.Before(time.Now())
}


type LRUCache struct {
	LCache
	mu sync.RWMutex
	list   	  *List
	liveTime  *time.Duration
	items     map[string]*Item
	keySize   int64			// ignore
	valueSize int64			// ignore
	limitSize int64			// list len
}

func newLRUCache(cfg *Config) *LRUCache {
	c := &LRUCache{
		list: newList(),
		items: map[string]*Item{},
		limitSize: cfg.size,
		liveTime: cfg.expire,
	}
	c.loadGroup = &Group{}
	c.StartClearTask()
	return c
}

func (c *LRUCache) Get(key string) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.get(key)
}

func (c *LRUCache) get(key string) interface{} {
	item, ok := c.items[key]
	if !ok {
		return nil
	}
	if item.isExpired() {
		delete(c.items, key)
		c.list.remove(item)
		return nil
	}
	c.list.transfer(item)
	return item.value
}

func (c *LRUCache) set(key string, value interface{}) *Item {
	var (
		item   *Item
		exists bool
	)
	if item, exists = c.items[key]; exists {
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
	return item
}


func (c *LRUCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item := c.set(key, value)
	if c.liveTime == nil {
		item.expire = nil
	} else {
		t := time.Now().Add(*c.liveTime)
		item.expire = &t
	}
}

func (c *LRUCache) SetWithExpire(key string, value interface{}, expire time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item := c.set(key, value)
	t := time.Now().Add(expire)
	item.expire = &t
}

func (c *LRUCache) Del(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, ok := c.items[key]; ok {
		c.list.remove(item)
		delete(c.items, key)
	}
}

func (c *LRUCache) Load(key string, loader func() (interface{}, error)) (interface{}, error) {
	c.mu.Lock()
	v := c.get(key)
	c.mu.Unlock()
	if v != nil {
		return v, nil
	}
	value, err := c.load(key, loader)
	return value, err
}

func (c *LRUCache) LoadNoWait(key string, loader func() (interface{}, error)) (interface{}, error) {
	return nil, nil
}

// loader: get data from source
func (c *LRUCache) load(key string, loader func() (interface{}, error)) (interface{}, error) {
	value, err := c.loadGroup.Do(key, func() (interface{}, error) {
		v, err := loader()
		if err != nil {
			return nil, err
		}
		c.Set(key, v)
		return v, nil
	})
	return value, err
}

func (c *LRUCache) Keys() []string {
	keys := make([]string, 0, len(c.items))
	c.mu.RLock()
	for k, _ := range c.items {
		keys = append(keys, k)
	}
	c.mu.RUnlock()
	return keys
}

// 这里并不严格删除过期的
func (c *LRUCache) clearExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	delNum := c.list.len
	if delNum == 0 {
		return
	}
	if delNum > 100 {
		delNum = 100
	}
	var nt *Item
	node := c.list.head.prev
	for i := 0; i < delNum; i++ {
		if !node.isExpired() {
			return
		}
		nt = node.prev
		c.list.remove(node)
		delete(c.items, node.key)
		node = nt
	}
}

func (c *LRUCache) StartClearTask() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <- ticker.C :
				c.clearExpired()
			}
		}
	}()
}
