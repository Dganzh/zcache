package zcache

import "time"

const (
	EvictLru = "lru"
)

var defaultConfig = Config{
	evictType: EvictLru,
	size: 100,
}


type ConfigOption func(config *Config)


type Config struct {
	evictType string
	size 	  int64
	expire  *time.Duration
}

func WithLru() ConfigOption {
	return func(cfg *Config) {
		cfg.evictType = EvictLru
	}
}

func WithSize(size int64) ConfigOption {
	return func(cfg *Config) {
		cfg.size = size
	}
}

func WithExpire(exp *time.Duration) ConfigOption {
	return func(cfg *Config) {
		cfg.expire = exp
	}
}

type DConfig struct {
	Config
	addr      string
	peerAddrs []string
}
