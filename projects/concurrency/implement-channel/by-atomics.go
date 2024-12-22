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
	value atomic.Pointer[string]
}

func (c *channel) Send(str string) {
	for {
		// if c.value.CompareAndSwap(nil, &str) {
		if c.value.Load() == nil {
			c.value.Store(&str)
			break
		}
	}
}

func (c *channel) Recv() string {
	for {
		if c.value.Load() != nil {
			// if c.value.CompareAndSwap(valueToReturn, nil) {
			valueToReturn := c.value.Load()
			c.value.Store(nil)
			return *valueToReturn
			// }
		}
	}
}
