package main

import (
	"fmt"
	"sync"
)

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

type static struct {
	readCount        int //Question 1
	unreadCount      int //Question 1
	entriesNeverRead int //Question 2
	totalReadExisted int //Question 3
	writesCount      int //Question 4
}
type Cache[K comparable, V any] struct {
	size int
	m    map[K]*Node[K, V]
	l    List[K, V]
	s    static
	rwm  sync.RWMutex
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

func NewStatic() *static {
	return &static{
		readCount:        0,
		unreadCount:      0,
		entriesNeverRead: 0,
		totalReadExisted: 0,
		writesCount:      0,
	}
}

func NewCache[K comparable, V any](entryLimit int) Cache[K, V] { //All K should be unique && same type && comparable type like primitive that can compare with each other

	return Cache[K, V]{
		size: entryLimit,
		m:    make(map[K]*Node[K, V]),
		l:    *NewList[K, V](),
		s:    *NewStatic(),
	}
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

func (c *Cache[K, V]) Put(key K, value V) bool {
	c.rwm.Lock()
	if currentNode, isExisted := c.m[key]; isExisted {
		currentNode.value = value
		c.l.moveNode(currentNode)
		c.rwm.Unlock()
		return true
	} else if len(c.m) >= c.size {
		deleted, entriesRead := c.l.deleteNode()
		delete(c.m, deleted)
		c.s.totalReadExisted -= entriesRead
	}
	newNode := c.l.addNode(key, value)
	c.m[key] = newNode
	c.s.entriesNeverRead++
	c.s.writesCount++
	c.rwm.Unlock()
	return false
}

func (c *Cache[K, V]) Get(key K) (*V, bool) {
	c.rwm.RLock()
	if currentNode, isExisted := c.m[key]; isExisted {
		c.l.moveNode(currentNode)
		c.s.readCount++
		c.s.totalReadExisted++
		currentNode.entriesRead++
		if currentNode.isRead == false {
			currentNode.isRead = true
			c.s.entriesNeverRead--
		}
		c.rwm.RUnlock()
		return &currentNode.value, true
	}
	c.s.unreadCount++
	c.rwm.RUnlock()
	return nil, false
}

var wg sync.WaitGroup

func main() {
	cache := NewCache[int, string](4)
	wg.Add(7)
	go func() {
		cache.Put(1, "newValue1")
		wg.Done()
	}()
	go func() {
		cache.Put(2, "newValue2")
		wg.Done()
	}()
	go func() {
		cache.Put(3, "newValue3")
		wg.Done()
	}()
	go func() {
		cache.Get(2)
		wg.Done()
	}()
	go func() {
		cache.Put(4, "newValue4")
		wg.Done()
	}()
	go func() {
		cache.Put(2, "newValue22")
		wg.Done()
	}()
	go func() {
		cache.Put(5, "newValue5")
		wg.Done()
	}()
	wg.Wait()
	wg.Add(1)
	go func() {
		if value, isExisted := cache.Get(2); isExisted {
			fmt.Printf("Value is:%v\n", *value)
		} else {
			fmt.Println("Value isn't exist")
		}
		wg.Done()
	}()
	wg.Wait()
	if totalHit := cache.s.readCount + cache.s.unreadCount; totalHit == 0 {
		fmt.Println("Never have had any read from this cache!")
	} else {
		fmt.Printf("Hit rate:%.2f\n", float64(cache.s.readCount)/float64(totalHit)*100)
	}
	fmt.Printf("Entries were written to the cache and have never been read: %d\n", cache.s.entriesNeverRead)
	fmt.Printf("Average number of times that things currently in the cache is read: %.2f\n", float64(cache.s.totalReadExisted)/float64(len(cache.m)))
	fmt.Printf("Total reads and writes have been performed in the cache including evicted: %d\n", cache.s.readCount+cache.s.writesCount)
}
