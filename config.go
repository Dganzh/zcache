package zcache


const (
	EvictLru = "lru"
)

var defaultConfig = Config{
	evictType: EvictLru,
	size: 3,
}


type ConfigOption func(config *Config)


type Config struct {
	evictType string
	size 	  int64
}

func WithLru() ConfigOption {
	return func(cfg *Config) {
		cfg.evictType = EvictLru
	}
}

type DConfig struct {
	Config
	addr      string
	peerAddrs []string
}
