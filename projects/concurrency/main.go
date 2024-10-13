package main

import (
	"concurrency/list"
	"fmt"
	"sync"
)

type static struct {
	readCount        int //Question 1
	unreadCount      int //Question 1
	entriesNeverRead int //Question 2
	totalReadExisted int //Question 3
	writesCount      int //Question 4
}
type Cache[K comparable, V any] struct {
	size int
	m    map[K]*list.Node[K, V]
	l    list.List[K, V]
	s    static
	mu   sync.Mutex
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
		m:    make(map[K]*list.Node[K, V]),
		l:    *list.NewList[K, V](),
		s:    *NewStatic(),
	}
}

func (c *Cache[K, V]) Put(key K, value V) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if currentNode, isExisted := c.m[key]; isExisted {
		currentNode.Value = value
		c.l.MoveNodeToLast(currentNode)
		return true
	} else if len(c.m) >= c.size {
		deleted, entriesRead := c.l.DeleteFirstNode()
		delete(c.m, deleted)
		c.s.totalReadExisted -= entriesRead
	}
	newNode := c.l.AddNode(key, value)
	c.m[key] = newNode
	c.s.entriesNeverRead++
	c.s.writesCount++
	return false
}

func (c *Cache[K, V]) Get(key K) (*V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if currentNode, isExisted := c.m[key]; isExisted {
		c.l.MoveNodeToLast(currentNode)
		c.s.readCount++
		c.s.totalReadExisted++
		currentNode.EntriesRead++
		if !currentNode.IsRead {
			currentNode.IsRead = true
			c.s.entriesNeverRead--
		}
		return &currentNode.Value, true
	}
	c.s.unreadCount++
	return nil, false
}

func hitRate(readCount int, unreadCount int) float32 {
	totalHit := readCount + unreadCount
	return float32(readCount) / float32(totalHit) * 100
}

func main() {
	var wg sync.WaitGroup
	cache := NewCache[int, string](4)
	wg.Add(3)
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
	wg.Wait()
	wg.Add(4)
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
	hitRate := hitRate(cache.s.readCount, cache.s.unreadCount)
	if hitRate == 0 {
		fmt.Println("Never have had any read from this cache!")
	} else {
		fmt.Printf("Hit rate:%.2f\n", hitRate)
	}
	fmt.Printf("Entries were written to the cache and have never been read: %d\n", cache.s.entriesNeverRead)
	fmt.Printf("Average number of times that things currently in the cache is read: %.2f\n", float32(cache.s.totalReadExisted)/float32(len(cache.m)))
	fmt.Printf("Total reads and writes have been performed in the cache including evicted: %d\n", cache.s.readCount+cache.s.writesCount)
}
