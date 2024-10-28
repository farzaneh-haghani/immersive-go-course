package computingCache

import (
	"concurrency/cache"
)

type computingCache[K comparable, V any] struct {
	cache   cache.Cache[K, V]
	Creator func(K) V
}

func NewComputingCache[K comparable, V any](entryLimit int, creator func(K) V) computingCache[K, V] {
	return computingCache[K, V]{
		cache:   cache.NewCache[K, V](entryLimit),
		Creator: creator,
	}
}

func Creator[K comparable](key int) any {
	value := "NewValue"
	return value
}

func (c *computingCache[K, V]) Get(key K) V {
	if currentValue, isExisted := c.cache.Get(key); isExisted {
		return *currentValue
	}
	value := c.Creator(key)
	c.cache.Put(key, value)
	return value
}
