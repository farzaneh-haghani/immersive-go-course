package computingCache

import (
	"concurrency/cache"
)

type computingCache[K comparable, V any] struct {
	cache.Cache[K, V]
	Creator func(K) V
}

func NewComputingCache[K comparable, V any](entryLimit int, creator func(K) V) computingCache[K, V] {
	return computingCache[K, V]{
		Cache:   cache.NewCache[K, V](entryLimit),
		Creator: creator,
	}
}

func Creator[K comparable](key int) any {
	value := "NewValue"
	return value
}

func (c *computingCache[K, V]) Get(key K) V {
	if currentNode, isExisted := c.M[key]; isExisted {
		return currentNode.Data.Value
	}
	if len(c.M) >= c.Size {
		deleted, entriesRead := c.L.DeleteFirstNode()
		delete(c.M, deleted)
		c.S.TotalReadExisted -= entriesRead
	}
	value := c.Creator(key)
	newNode := c.L.AddNodeToLast(key, value)
	c.M[key] = newNode
	return value
}
