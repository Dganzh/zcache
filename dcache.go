package zcache

/*
distributed cache
分布式缓存缓存
*/

import (
	"github.com/Dganzh/zrpc"
	"strings"
	"time"
)

type RpcDCache struct {
	cache   Cache
	server  *zrpc.Server
	peers   []*zrpc.Client
	sendChs []chan *string
	recvCh  chan string
}


func NewDCache(addr string, peers []string, opts ...ConfigOption) DCache {
	clients := make([]*zrpc.Client, len(peers))
	chs := make([]chan *string, len(peers))
	for i := 0; i < len(peers); i++ {
		clients[i] = zrpc.NewClient(peers[i])
		chs[i] = make(chan *string, 1000)
	}
	server := zrpc.NewServer(addr)
	r := &RpcDCache{
		server: server,
		peers: clients,
		sendChs: chs,
		cache: NewCache(opts...),
	}
	r.sendSyncReq()
	server.Register(r)
	go server.Start()
	return r
}

func (c *RpcDCache) Exit() {
	//c.server.Stop()
}

func (c *RpcDCache) Get(key string) interface{} {
	return c.cache.Get(key)
}

func (c *RpcDCache) Set(key string, value interface{}) {
	c.cache.Set(key, value)
}

func (c *RpcDCache) Del(key string) {
	c.cache.Del(key)
	c.SyncDel(key)
}

func (c *RpcDCache) Load(key string, loader func() (interface{}, error)) (interface{}, error) {
	return c.cache.Load(key, loader)
}

func (c *RpcDCache) SyncDel(key string) {
	for i := 0; i < len(c.sendChs); i++ {
		c.sendChs[i] <- &key
	}
}

func (c *RpcDCache) sendSyncReq() {
	for idx := 0; idx < len(c.sendChs); idx++ {
		go c.sendPeerSyncReq(idx)
	}
}

func (c *RpcDCache) sendPeerSyncReq(idx int) {
	sendSize := 100
	retryTimes := 3
	shouldSend := false
	toSend := make([]string, 0, sendSize)
	sendCh := c.sendChs[idx]
	peer := c.peers[idx]
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case key := <- sendCh :
			toSend = append(toSend, *key)
			shouldSend = len(toSend) >= sendSize
		case <- ticker.C:
			shouldSend = len(toSend) > 0
		}
		if shouldSend {
			keys := strings.Join(toSend, "$$")
			for i := 0; i < retryTimes; i++ {
				if peer.Call("RpcDCache.SyncHandler", &keys, nil) {
					break
				}
			}
			toSend = toSend[:0]
			//shouldSend = false
		}
	}
}

func (c *RpcDCache) SyncHandler(keys *string, reply *interface{}) {
	for _, key := range strings.Split(*keys, "$$") {
		c.cache.Del(key)
	}
}
