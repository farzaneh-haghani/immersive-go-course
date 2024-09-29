package main

import "fmt"

type Value[K comparable] struct {
	data any
	next *K
	prev *K
}

type Cache[K comparable] struct {
	size  int
	m     map[K]Value[K]
	first *K
	last  *K
}

func NewCache[K comparable](entryLimit int) Cache[K] { //All K should be unique && same type && comparable type like primitive that can compare with each other
	return Cache[K]{
		size:  entryLimit,
		m:     make(map[K]Value[K]),
		first: nil,
		last:  nil,
	}
}

func (c *Cache[K]) Put(key K, value Value[K]) bool {

	if currentValue, isExisted := c.m[key]; isExisted {
		nextKey := *currentValue.next
		nextValue := c.m[nextKey]

		if *c.first == key {
			nextValue.prev = nil
			*c.first = nextKey
		} else {
			nextValue.prev = currentValue.prev
			prevKey := *currentValue.prev
			prevValue := c.m[prevKey]
			prevValue.next = currentValue.next
			c.m[prevKey] = prevValue
		}
		c.m[nextKey] = nextValue
		lastValue := c.m[*c.last]
		lastValue.next = &key
		c.m[*c.last] = lastValue
		currentValue.prev = c.last
		currentValue.next = nil
		currentValue.data = value.data
		c.m[key] = currentValue
		c.last = &key
		return true
	} else if len(c.m) < c.size {

		if c.first == nil {
			value.prev = nil
			c.first = &key
		} else {
			oldLastValue := c.m[*c.last]
			oldLastValue.next = &key
			c.m[*c.last] = oldLastValue
			value.prev = c.last
		}
		value.next = nil
		c.last = &key
		c.m[key] = value
		return false
	} else {
		currentLastValue := c.m[*c.last]
		currentLastValue.next = &key
		c.m[*c.last] = currentLastValue

		currentFirstValue := c.m[*c.first]
		delete(c.m, *c.first)
		c.first = currentFirstValue.next

		value.prev = c.last
		value.next = nil
		c.m[key] = value

		c.last = &key
		return false
	}
}

func (c *Cache[K]) Get(key K) (*Value[K], bool) {
	if value, isExisted := c.m[key]; isExisted {
		return &value, true
	}
	return nil, false
}

func main() {
	cache := NewCache[int](5)
	newValue1 := Value[int]{data: "test"}
	newValue2 := Value[int]{data: "test2"}
	newValue3 := Value[int]{data: "test3"}
	cache.Put(1, newValue1)
	cache.Put(2, newValue2)
	cache.Put(3, newValue1)
	cache.Put(4, newValue1)
	cache.Put(5, newValue2)
	cache.Put(6, newValue3)
	cache.Put(7, newValue3)
	value, isExisted := cache.Get(7)
	if isExisted {
		fmt.Printf("Value is:%s", value.data)
	} else {
		fmt.Println("Value isn't exist")
	}
}
