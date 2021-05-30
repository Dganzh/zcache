package zcache


// 本地缓存的接口
type Cache interface {
	Set(string, interface{})
	Get(string) interface{}
	Del(string)
}


// 分布式缓存的接口
type DCache interface {
	Cache
	Syncer
}


// 同步其他节点接口
// 目前，只考虑缓存失效才广播
type Syncer interface {
	SyncDel(string)
	SyncHandler(*string, *interface{})
}

