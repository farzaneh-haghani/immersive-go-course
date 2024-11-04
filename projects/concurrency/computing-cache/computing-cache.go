package computingCache

import (
	"concurrency/cache"
	"sync"
	"time"
)

type computingCache[K comparable, V any] struct {
	cache           cache.Cache[K, V]
	Creator         func(K) V
	pendingRequests map[K]time.Time
	responseTime    time.Duration
	mu              sync.Mutex
}

func NewComputingCache[K comparable, V any](entryLimit int, Creator func(K) V) computingCache[K, V] {
	return computingCache[K, V]{
		cache:           cache.NewCache[K, V](entryLimit),
		Creator:         Creator,
		pendingRequests: map[K]time.Time{},
		responseTime:    0,
	}
}

func Creator[K comparable](key int) string {
	value := "NewComputingCacheValue"
	return value
}

func (c *computingCache[K, V]) SendRequest(key K) V {
	c.mu.Lock()
	c.pendingRequests[key] = time.Now()
	c.mu.Unlock()
	value := c.Creator(key)
	c.cache.Put(key, value)
	c.mu.Lock()
	delete(c.pendingRequests, key)
	c.mu.Unlock()
	return value
}

func (c *computingCache[K, V]) Get(key K) V {
	if currentValue, isExisted := c.cache.Get(key); isExisted {
		return *currentValue
	}
	c.mu.Lock()
	if requestedTime, isExisted := c.pendingRequests[key]; isExisted {
		c.mu.Unlock()
		for {
			if time.Since(requestedTime) > c.responseTime+1000000000 {
				value := c.SendRequest(key)
				return value
			}
			time.Sleep(time.Since(requestedTime) + 1000000000)
			if value, isExisted := c.cache.Get(key); isExisted {
				c.mu.Lock()
				delete(c.pendingRequests, key)
				c.mu.Unlock()
				return *value
			}
		}
	}
	c.mu.Unlock()
	value := c.SendRequest(key)
	return value
}
