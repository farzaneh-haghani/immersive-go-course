package customCache

import (
	"concurrency/cache"
	"sync"
)

type CustomValue[K comparable] struct {
	data        any
	next        *K
	prev        *K
	IsRead      bool
	EntriesRead int
}

type CustomCache[K comparable] struct {
	size  int
	m     map[K]CustomValue[K]
	first *K
	last  *K
	S     cache.Static
	mu    sync.Mutex
}

func NewCustomValue[K comparable](value any) CustomValue[K] {
	return CustomValue[K]{
		data:        value,
		next:        nil,
		prev:        nil,
		IsRead:      false,
		EntriesRead: 0,
	}
}

func NewCustomCache[K comparable](entryLimit int) CustomCache[K] {
	return CustomCache[K]{
		size:  entryLimit,
		m:     make(map[K]CustomValue[K]),
		first: nil,
		last:  nil,
		S:     *cache.NewStatic(),
	}
}

func (c *CustomCache[K]) moveValueToLast(key K) {
	if *c.last == key {
		return
	}
	currentValue := c.m[key]

	lastValue := c.m[*c.last]
	lastValue.next = &key
	c.m[*c.last] = lastValue

	if *c.first == key {
		c.first = currentValue.next
		firstValue := c.m[*c.first]
		firstValue.prev = nil
		c.m[*c.first] = firstValue

	} else {
		prevKey := *currentValue.prev
		prevValue := c.m[prevKey]

		nextKey := *currentValue.next
		nextValue := c.m[nextKey]

		prevValue.next = &nextKey
		nextValue.prev = &prevKey

		c.m[prevKey] = prevValue
		c.m[nextKey] = nextValue
	}

	currentValue.prev = c.last
	c.last = &key
	currentValue.next = nil
	c.m[key] = currentValue
}

func (c *CustomCache[K]) deleteFirstValue() int {
	oldFirstKey := *c.first
	entriesRead := c.m[*c.first].EntriesRead
	if *c.first == *c.last {
		c.first = nil
		c.last = nil
	} else {
		*c.first = *c.m[*c.first].next
		firstValue := c.m[*c.first]
		firstValue.prev = nil
		c.m[*c.first] = firstValue
	}
	delete(c.m, oldFirstKey)
	return entriesRead
}

func (c *CustomCache[K]) addValue(key K, value any) {
	newValue := NewCustomValue[K](value)

	if c.first == nil {
		c.first = &key
	} else {
		lastValue := c.m[*c.last]
		lastValue.next = &key
		c.m[*c.last] = lastValue

		newValue.prev = c.last
	}
	c.m[key] = newValue
	c.last = &key
}

func (c *CustomCache[K]) Put(key K, value any) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if currentValue, isExisted := c.m[key]; isExisted {
		currentValue.data = value
		c.m[key] = currentValue

		c.moveValueToLast(key)
		return true
	} else if len(c.m) >= c.size {
		entriesRead := c.deleteFirstValue()
		c.S.TotalReadExisted -= entriesRead
	}
	c.addValue(key, value)
	c.S.EntriesNeverRead++
	c.S.WritesCount++
	return false
}

func (c *CustomCache[K]) Get(key K) (*any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if value, isExisted := c.m[key]; isExisted {
		c.moveValueToLast(key)
		c.S.ReadCount++
		c.S.TotalReadExisted++
		value.EntriesRead++
		if !value.IsRead {
			value.IsRead = true
			c.m[key] = value
			c.S.EntriesNeverRead--
		}
		return &value.data, true
	}
	c.S.UnreadCount++
	return nil, false
}
