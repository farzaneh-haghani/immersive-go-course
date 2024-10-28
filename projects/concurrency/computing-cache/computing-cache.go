package computingCache

import (
	"concurrency/cache"
	"time"
)

type computingCache[K comparable, V any] struct {
	cache       cache.Cache[K, V]
	Creator     func(K) V
	isRequested bool
}

func NewComputingCache[K comparable, V any](entryLimit int, creator func(K) V) computingCache[K, V] {
	return computingCache[K, V]{
		cache:       cache.NewCache[K, V](entryLimit),
		Creator:     creator,
		isRequested: false,
	}
}

func Creator[K comparable](key int) string {
	value := "NewComputingCacheValue"
	return value
}

func (c *computingCache[K, V]) Get(key K) V {
	if currentValue, isExisted := c.cache.Get(key); isExisted {
		return *currentValue
	}
	if c.isRequested {
		for {
			time.Sleep(time.Second)
			if currentValue, isExisted := c.cache.Get(key); isExisted {
				return *currentValue
			}
		}
	} else {
		c.isRequested = true
		value := c.Creator(key)
		c.cache.Put(key, value)
		return value
	}
}
