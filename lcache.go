package zcache
/*
local cache
实现基本的本地缓存
 */


type LCache struct {
	cache Cache
	config *Config

}


func NewCache(opts ...ConfigOption) Cache {
	cfg := defaultConfig
	for _, o := range opts {
		o(&cfg)
	}
	var cache Cache
	switch cfg.evictType {
	case EvictLru:
		cache = newLRUCache(&cfg)
	default:
		panic("zcache: Unknown evict type " + cfg.evictType)
	}
	return cache
}

