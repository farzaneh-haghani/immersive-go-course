package main

import "fmt"

type Cache[K comparable, V any] struct {
	size    int
	m       map[K]V
}

func NewCache[K comparable, V any](entryLimit int) Cache[K, V] {
	return Cache[K, V]{
		size: entryLimit,
		m:    make(map[K]V),
	}
}

func (c *Cache[K, V]) Put(key K, value V) bool {
	if _, isExisted := c.m[key]; isExisted {
		c.m[key] = value
		return true
	} else if len(c.m) < c.size {
		c.m[key] = value
		return false
	} else {

		return false
	}

}

func (c *Cache[K, V]) Get(key K) (*V, bool) {
	if value, isExisted := c.m[key]; isExisted {
		return &value, true
	}
	return nil, false
}

func main() {
	cache := NewCache[int, string](5)
	cache.Put(1, "test")
	cache.Put(1, "test2")
	value, _ := cache.Get(1)
	fmt.Printf("Value is:%s", *value)

}
