package cache

import (
	"github.com/cultureamp/glamplify/helper"
	cachego "github.com/patrickmn/go-cache"
	"time"
)

// Config for the cache
type Config struct {
	CacheDuration time.Duration
}

// Cache represents a simple cache
type Cache struct {
	conf Config
	cache *cachego.Cache
}

// New creates a new Cache
func New(configure ...func(*Config)) *Cache {

	c := helper.GetEnvInt(CacheDurationEnv, 60)
	cacheDuration := time.Duration(c) * time.Second

	conf := Config{
		CacheDuration:  cacheDuration,
	}

	for _, config := range configure {
		config(&conf)
	}

	cache := cachego.New(conf.CacheDuration, conf.CacheDuration)

	return &Cache {
		conf: conf,
		cache: cache,
	}
}

// Get an item from the cache
func (c Cache) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

// Set an item within the cache
func (c *Cache) Set(key string, val interface{}, duration time.Duration) {
	c.cache.Set(key, val, duration)
}

