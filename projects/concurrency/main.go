package main

import (
	"concurrency/cache"
	customCache "concurrency/custom-cache"
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	length := 4
	cache := cache.NewCache[int, string](length)
	customCache := customCache.NewCustomCache[int](length)

	wg.Add(3)
	go func() {
		cache.Put(1, "newValue1")
		customCache.Put(1, "newValue1")
		wg.Done()
	}()
	go func() {
		cache.Put(2, "newValue2")
		customCache.Put(2, "newValue2")
		wg.Done()
	}()
	go func() {
		cache.Put(3, "newValue3")
		customCache.Put(3, "newValue3")
		wg.Done()
	}()
	wg.Wait()
	wg.Add(4)
	go func() {
		cache.Get(2)
		customCache.Get(2)
		wg.Done()
	}()
	go func() {
		cache.Put(4, "newValue4")
		customCache.Put(4, "newValue4")
		wg.Done()
	}()
	go func() {
		cache.Put(2, "newValue22")
		customCache.Put(2, "newValue22")
		wg.Done()
	}()
	go func() {
		cache.Put(5, "newValue5")
		customCache.Put(5, "newValue5")
		wg.Done()
	}()
	wg.Wait()
	wg.Add(1)
	go func() {
		if value, isExisted := cache.Get(2); isExisted {
			fmt.Printf("Value of cache is: %v\n", *value)
		} else {
			fmt.Println("Value of cache doesn't exist")
		}

		if customCache, isExisted := customCache.Get(2); isExisted {
			fmt.Printf("Value of custom cache is: %v\n", *customCache)
		} else {
			fmt.Println("Value doesn't exist")
		}
		wg.Done()
	}()
	wg.Wait()

	fmt.Println("Statics for cache:")
	cache.PrintStatics(cache.S, length)
	fmt.Println("Statics for custom cache:")
	cache.PrintStatics(customCache.S, length)
}
