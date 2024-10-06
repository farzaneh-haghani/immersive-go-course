package main

import "fmt"

type Node[K comparable, V any] struct {
	key         K
	value       V
	next        *Node[K, V]
	prev        *Node[K, V]
	isRead      bool //Question 2
	entriesRead int  //Question 3
}

type List[K comparable, V any] struct {
	first *Node[K, V]
	last  *Node[K, V]
}

type Cache[K comparable, V any] struct {
	size             int
	m                map[K]*Node[K, V]
	readCount        int //Question 1
	unreadCount      int //Question 1
	entriesNeverRead int //Question 2
	totalReadExisted int //Question 3
	writesCount      int //Question 4
}

func NewNode[K comparable, V any](key K, value V) *Node[K, V] {
	return &Node[K, V]{
		key:         key,
		value:       value,
		next:        nil,
		prev:        nil,
		isRead:      false,
		entriesRead: 0,
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
		size:             entryLimit,
		m:                make(map[K]*Node[K, V]),
		readCount:        0,
		unreadCount:      0,
		entriesNeverRead: 0,
		totalReadExisted: 0,
		writesCount:      0,
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

func (l *List[K, V]) deleteNode() (K, int) {
	deleted := l.first.key
	entriesRead := l.first.entriesRead
	if l.first == l.last {
		l.first = nil
		l.last = nil
	} else {
		l.first = l.first.next
		l.first.prev = nil
	}
	return deleted, entriesRead
}

func (c *Cache[K, V]) Put(key K, l *List[K, V], value V) bool {
	if currentNode, isExisted := c.m[key]; isExisted {
		currentNode.value = value
		l.moveNode(currentNode)
		return true
	} else if len(c.m) >= c.size {
		deleted, entriesRead := l.deleteNode()
		delete(c.m, deleted)
		c.totalReadExisted -= entriesRead
	}
	newNode := l.addNode(key, value)
	c.m[key] = newNode
	c.entriesNeverRead++
	c.writesCount++
	return false
}

func (c *Cache[K, V]) Get(key K, l *List[K, V]) (*Node[K, V], bool) {
	if currentNode, isExisted := c.m[key]; isExisted {
		l.moveNode(currentNode)
		c.readCount++
		c.totalReadExisted++
		currentNode.entriesRead++
		if currentNode.isRead == false {
			currentNode.isRead = true
			c.entriesNeverRead--
		}
		return currentNode, true
	}
	c.unreadCount++
	return nil, false
}

func main() {
	cache, list := NewCache[int, string](4)
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
	fmt.Printf("Entries were written to the cache and have never been read: %d\n", cache.entriesNeverRead)
	fmt.Printf("Average number of times that things currently in the cache is read: %.2f\n", float64(cache.totalReadExisted)/float64(len(cache.m)))
	fmt.Printf("Total reads and writes have been performed in the cache including evicted: %d\n", cache.readCount+cache.writesCount)
}
