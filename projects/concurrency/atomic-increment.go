package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// var x = 0
var x atomic.Int32

func increment(wg *sync.WaitGroup) {
	// x = x + 1
	x.Add(1)
	wg.Done()
}

func main() {
	for i := 0; i < 5; i++ {
		var w sync.WaitGroup
		for i := 0; i < 1000; i++ {
			w.Add(1)
			go increment(&w)
		}
		w.Wait()
		// fmt.Println("final value of x", x)
		fmt.Println("final value of x", x.Load())
	}
}

