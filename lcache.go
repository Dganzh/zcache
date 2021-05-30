package zcache
/*
local cache
实现基本的本地缓存
 */


type LCache struct {
	cache Cache
	config *Config

}


func NewLCache(opts ...ConfigOption) *LCache {
	cfg := &Config{}
	for _, o := range opts {
		o(cfg)
	}
	var cache Cache
	switch cfg.evictType {
	case EvictLru:
		cache = newLRUCache()
	default:
		panic("zcache: Unknown evict type " + cfg.evictType)
	}
	return &LCache{
		cache: cache,
		config: cfg,
	}
}

