//  * Use atomics instead of mutexes

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

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
	return &channel{}
}

type channel struct {
	isSent atomic.Bool
	value  atomic.Pointer[string]
}

func (c *channel) Send(str string) {
	for {
		if !c.isSent.Load() {
			c.value.Store(&str)
			c.isSent.Store(true)
			break
		}
	}
}

func (c *channel) Recv() string {
	for {
		if c.isSent.Load() {
			valueToReturn := *c.value.Load()
			c.isSent.Store(false)
			return valueToReturn
		}
	}
}
