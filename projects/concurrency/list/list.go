package list

type Entry[K comparable, V any] struct {
	key         K
	Value       V
	IsRead      bool //Question 2
	EntriesRead int  //Question 3
}
type Node[K comparable, V any] struct {
	Data Entry[K, V]
	next *Node[K, V]
	prev *Node[K, V]
}

type List[K comparable, V any] struct {
	first *Node[K, V]
	last  *Node[K, V]
}

func NewEntry[K comparable, V any](key K, value V) Entry[K, V] {
	return Entry[K, V]{
		key:         key,
		Value:       value,
		IsRead:      false,
		EntriesRead: 0,
	}
}

func NewNode[K comparable, V any](key K, value V) *Node[K, V] {
	return &Node[K, V]{
		Data: NewEntry(key, value),
		next: nil,
		prev: nil,
	}
}

func NewList[K comparable, V any]() *List[K, V] {
	return &List[K, V]{
		first: nil,
		last:  nil,
	}
}

func (l *List[K, V]) AddNodeToLast(key K, value V) *Node[K, V] {
	newNode := NewNode[K, V](key, value)
	if l.first == nil {
		l.first = newNode
	} else {
		l.last.next = newNode
		newNode.prev = l.last
	}
	l.last = newNode
	return newNode
}

func (l *List[K, V]) MoveNodeToLast(currentNode *Node[K, V]) {
	if l.last == currentNode {
		return
	}
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

func (l *List[K, V]) DeleteFirstNode() (K, int) {
	deleted := l.first.Data.key
	entriesRead := l.first.Data.EntriesRead
	if l.first == l.last {
		l.first = nil
		l.last = nil
	} else {
		l.first = l.first.next
		l.first.prev = nil
	}
	return deleted, entriesRead
}
