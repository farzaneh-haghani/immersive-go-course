package main

import "fmt"

type Node[K comparable, V any] struct {
	key              K
	value            V
	next             *Node[K, V]
	prev             *Node[K, V]
	entriesNeverRead bool
}

type List[K comparable, V any] struct {
	first *Node[K, V]
	last  *Node[K, V]
}

type Cache[K comparable, V any] struct {
	size        int
	m           map[K]*Node[K, V]
	readCount   int
	unreadCount int
	entries     int
}

func NewNode[K comparable, V any](key K, value V) *Node[K, V] {
	return &Node[K, V]{
		key:              key,
		value:            value,
		next:             nil,
		prev:             nil,
		entriesNeverRead: false,
	}
}

func NewList[K comparable, V any]() *List[K, V] {
	return &List[K, V]{
		first: nil,
		last:  nil,
	}
}

func NewCache[K comparable, V any](entryLimit int) (Cache[K, V], *List[K, V]) { //All K should be unique && same type && comparable type like primitive that can compare with each other
	list := NewList[K, V]()
	return Cache[K, V]{
		size:        entryLimit,
		m:           make(map[K]*Node[K, V]),
		readCount:   0,
		unreadCount: 0,
		entries:     0,
	}, list
}

func (l *List[K, V]) addNode(key K, value V) *Node[K, V] {
	newNode := NewNode(key, value)
	if l.first == nil {
		l.first = newNode
	} else {
		l.last.next = newNode
		newNode.prev = l.last
	}
	l.last = newNode
	return newNode
}

func (l *List[K, V]) moveNode(currentNode *Node[K, V]) {
	if l.last != currentNode {
		l.last.next = currentNode
		if l.first != currentNode {
			prevNode := currentNode.prev
			prevNode.next = currentNode.next
			nextNode := currentNode.next
			nextNode.prev = currentNode.prev
		} else {
			l.first = currentNode.next
			l.first.prev = nil
		}
		currentNode.prev = l.last
		l.last = currentNode
		currentNode.next = nil
	}
}

func (l *List[K, V]) deleteNode() K {
	deleted := l.first.key
	if l.first == l.last {
		l.first = nil
		l.last = nil
	} else {
		l.first = l.first.next
		l.first.prev = nil
	}
	return deleted
}

func (c *Cache[K, V]) Put(key K, l *List[K, V], value V) bool {

	if currentNode, isExisted := c.m[key]; isExisted {
		currentNode.value = value
		l.moveNode(currentNode)
		return true
	} else if len(c.m) >= c.size {
		deleted := l.deleteNode()
		delete(c.m, deleted)
	}
	newNode := l.addNode(key, value)
	c.m[key] = newNode
	c.entries++
	return false
}

func (c *Cache[K, V]) Get(key K, l *List[K, V]) (*Node[K, V], bool) {
	if currentNode, isExisted := c.m[key]; isExisted {
		l.moveNode(currentNode)
		c.readCount++
		if currentNode.entriesNeverRead == false {
			currentNode.entriesNeverRead = true
			c.entries--
		}
		return currentNode, true
	}
	c.unreadCount++
	return nil, false
}

func main() {
	cache, list := NewCache[int, string](2)
	cache.Put(1, list, "newValue1")
	cache.Put(2, list, "newValue2")
	cache.Put(3, list, "newValue3")
	cache.Get(2, list)
	cache.Put(4, list, "newValue4")
	cache.Put(2, list, "newValue22")
	cache.Put(5, list, "newValue5")
	if value, isExisted := cache.Get(2, list); isExisted {
		fmt.Printf("Value is:%s\n", value.value)
	} else {
		fmt.Println("Value isn't exist\n")
	}
	if totalHit := cache.readCount + cache.unreadCount; totalHit == 0 {
		fmt.Println("Never have had any read from this cache!")
	} else {
		fmt.Printf("Hit rate:%d\n", cache.readCount/totalHit*100)
	}
	fmt.Printf("Entries were written to the cache and have never been read: %d\n", cache.entries)
}
