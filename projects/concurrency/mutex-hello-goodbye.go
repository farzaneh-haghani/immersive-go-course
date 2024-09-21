package main

import (
	"fmt"
	"sync"
)

// Tasks to try out:
//  * Use atomics instead of mutexes
//  * Write a buffered channel (with known capacity up-front)
//  * Discuss the trade-offs of a buffered channel with fixed capacity vs a dynamically growable buffered channel
//  * Write some tests showing why the locks are needed (i.e. tests that fail if we don't do all of the locking)
//  * Try out different APIs (e.g. Recv with timeout, Recv which returns nil if there's no value)

func main() {
	ch := NewUnbufferedChannel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		ch.Send("hello")
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		ch.Send("goodbye")
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		value := ch.Recv()
		fmt.Println(value)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		value := ch.Recv()
		fmt.Println(value)
		wg.Done()
	}()

	wg.Wait()
}

func NewUnbufferedChannel() *channel {
	return &channel{
		value: nil,
	}
}

type channel struct {
	mu    sync.Mutex
	value *string
}

func (c *channel) Send(str string) {
	// Set value to the string being sent
	for {
		c.mu.Lock()
		if c.value == nil {
			c.value = &str
			c.mu.Unlock()
			break
		}
		c.mu.Unlock()
	}
}

func (c *channel) Recv() string {
	// Wait until value isn't nil
	for {
		c.mu.Lock()
		if c.value != nil {
			valueToReturn := *c.value
			c.value = nil
			c.mu.Unlock()
			return valueToReturn
		}
		c.mu.Unlock()
	}
}

