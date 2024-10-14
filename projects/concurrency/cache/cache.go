package cache

import (
	"concurrency/list"
	"fmt"
	"sync"
)

type Static struct {
	ReadCount        int //Question 1
	UnreadCount      int //Question 1
	EntriesNeverRead int //Question 2
	TotalReadExisted int //Question 3
	WritesCount      int //Question 4
}

type Cache[K comparable, V any] struct {
	Size int
	m    map[K]*list.Node[K, V]
	l    list.List[K, V]
	S    Static
	mu   sync.Mutex
}

func NewStatic() *Static {
	return &Static{
		ReadCount:        0,
		UnreadCount:      0,
		EntriesNeverRead: 0,
		TotalReadExisted: 0,
		WritesCount:      0,
	}
}

func NewCache[K comparable, V any](entryLimit int) Cache[K, V] { //All K should be unique && same type && comparable type like primitive that can compare with each other
	return Cache[K, V]{
		Size: entryLimit,
		m:    make(map[K]*list.Node[K, V]),
		l:    *list.NewList[K, V](),
		S:    *NewStatic(),
	}
}

func (c *Cache[K, V]) Put(key K, value V) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if currentNode, isExisted := c.m[key]; isExisted {
		currentNode.Value = value
		c.l.MoveNodeToLast(currentNode)
		return true
	} else if len(c.m) >= c.Size {
		deleted, entriesRead := c.l.DeleteFirstNode()
		delete(c.m, deleted)
		c.S.TotalReadExisted -= entriesRead
	}
	newNode := c.l.AddNodeToLast(key, value)
	c.m[key] = newNode
	c.S.EntriesNeverRead++
	c.S.WritesCount++
	return false
}

func (c *Cache[K, V]) Get(key K) (*V, bool) {
	// c.rwm.RLock() is not safe because we are not just reading, we move the node as well, Also RWMutex is good when we use read more than write to make it faster. but when we write more than read, it makes it slower because it always should check RLock extra.
	c.mu.Lock()
	defer c.mu.Unlock()
	if currentNode, isExisted := c.m[key]; isExisted {
		c.l.MoveNodeToLast(currentNode)
		c.S.ReadCount++
		c.S.TotalReadExisted++
		currentNode.EntriesRead++
		if !currentNode.IsRead {
			currentNode.IsRead = true
			c.S.EntriesNeverRead--
		}
		return &currentNode.Value, true
	}
	c.S.UnreadCount++
	return nil, false
}

func (c *Cache[K, V]) PrintStatics(static Static, length int) {
	totalHit := static.ReadCount + static.UnreadCount
	hitRate := float32(static.ReadCount) / float32(totalHit) * 100
	if hitRate == 0 {
		fmt.Println("Never have had any read from this cache!")
	} else {
		fmt.Printf("Hit rate:%.2f\n", hitRate)
	}
	fmt.Printf("Entries were written to the cache and have never been read: %d\n", static.EntriesNeverRead)
	fmt.Printf("Average number of times that things currently in the cache is read: %.2f\n", float32(static.TotalReadExisted)/float32(length))
	fmt.Printf("Total reads and writes have been performed in the cache including evicted: %d\n", static.ReadCount+static.WritesCount)

}
