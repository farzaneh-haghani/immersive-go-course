package computingCache

import (
	"concurrency/cache"
	"sync"
)

type computingCache[K comparable, V any] struct {
	cache           cache.Cache[K, V]
	Creator         func(K) V
	pendingRequests map[K][]chan V // My brain didn't work more so I cheated the idea of slice channel from your sample solutions
	// responseTime    time.Duration
	mu sync.Mutex
}

func NewComputingCache[K comparable, V any](entryLimit int, Creator func(K) V) computingCache[K, V] {
	return computingCache[K, V]{
		cache:           cache.NewCache[K, V](entryLimit),
		Creator:         Creator,
		pendingRequests: make(map[K][]chan V),
		// responseTime:    0,
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
	c.mu.Lock()
	if channels, isExisted := c.pendingRequests[key]; isExisted {
		newChannel := make(chan V)
		c.pendingRequests[key] = append(channels, newChannel)
		c.mu.Unlock()

		value := <-newChannel //If the first go routine after writing in the map, lose connection, all the rest go routines will stay forever?
		return value
	}
	c.pendingRequests[key] = []chan V{}
	c.mu.Unlock()
	value := c.Creator(key)
	c.cache.Put(key, value)
	c.mu.Lock()
	for _, channel := range c.pendingRequests[key] {
		channel <- value
	}
	delete(c.pendingRequests, key)
	c.mu.Unlock()
	return value
}
