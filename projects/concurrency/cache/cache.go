package cache

import (
	"concurrency/list"
	"fmt"
	"io"
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
	M    map[K]*list.Node[K, V]
	L    list.List[K, V]
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
		M:    make(map[K]*list.Node[K, V]),
		L:    *list.NewList[K, V](),
		S:    *NewStatic(),
	}
}

func (c *Cache[K, V]) Put(key K, value V) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if currentNode, isExisted := c.M[key]; isExisted {
		currentNode.Data.Value = value
		c.L.MoveNodeToLast(currentNode)
		return true
	}
	if len(c.M) >= c.Size {
		deleted, entriesRead := c.L.DeleteFirstNode()
		delete(c.M, deleted)
		c.S.TotalReadExisted -= entriesRead
	}
	newNode := c.L.AddNodeToLast(key, value)
	c.M[key] = newNode
	c.S.EntriesNeverRead++
	c.S.WritesCount++
	return false
}

func (c *Cache[K, V]) Get(key K) (*V, bool) {
	// c.rwm.RLock() is not safe because we are not just reading, we move the node as well, Also RWMutex is good when we use read more than write to make it faster. but when we write more than read, it makes it slower because it always should check RLock extra.
	c.mu.Lock()
	defer c.mu.Unlock()
	if currentNode, isExisted := c.M[key]; isExisted {
		c.L.MoveNodeToLast(currentNode)
		c.S.ReadCount++
		c.S.TotalReadExisted++
		currentNode.Data.EntriesRead++
		if !currentNode.Data.IsRead {
			currentNode.Data.IsRead = true
			c.S.EntriesNeverRead--
		}
		return &currentNode.Data.Value, true
	}
	c.S.UnreadCount++
	return nil, false
}

func (c *Cache[K, V]) PrintStatics(w io.Writer, static Static, length int) {
	totalHit := static.ReadCount + static.UnreadCount
	hitRate := float32(static.ReadCount) / float32(totalHit) * 100
	if hitRate == 0 {
		io.WriteString(w, "Never have had any read from this cache!")
	} else {
		hit := fmt.Sprintf("Hit rate: %.2f\n", hitRate)
		io.WriteString(w, hit)
	}

	entries := fmt.Sprintf("Entries were written to the cache and have never been read: %d\n", static.EntriesNeverRead)
	io.WriteString(w, entries)
	avg := fmt.Sprintf("Average number of times that things currently in the cache is read: %.2f\n", float32(static.TotalReadExisted)/float32(length))
	io.WriteString(w, avg)
	total := fmt.Sprintf("Total reads and writes have been performed in the cache including evicted: %d\n", static.ReadCount+static.WritesCount)
	io.WriteString(w, total)
}
