package cache

import (
	"github.com/cultureamp/glamplify/helper"
	cachego "github.com/patrickmn/go-cache"
	"time"
)

type Config struct {
	CacheDuration time.Duration
}

type Cache struct {
	conf Config
	cache *cachego.Cache
}

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

func (c Cache) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

func (c *Cache) Set(key string, val interface{}, duration time.Duration) {
	c.cache.Set(key, val, duration)
}

