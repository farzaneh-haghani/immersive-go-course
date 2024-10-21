package main

import (
	"bytes"
	"concurrency/cache"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {

	t.Run("Test as a thread safe cache", func(t *testing.T) {
		var wg sync.WaitGroup
		length := 3
		cache := cache.NewCache[int, int](length) // Create a cache

		for i := 1; i < 4; i++ {
			wg.Add(1)
			go func() {
				cache.Put(i, i+5) // Put in the cache by 3 go routines
				wg.Done()
			}()
		}

		wg.Wait()
		wg.Add(2)
		go func() {
			value, isExisted := cache.Get(1) // Get from the cache
			require.Equal(t, *value, 6)
			require.Equal(t, isExisted, true)
			wg.Done()
		}()

		go func() {
			value, isExisted := cache.Get(2)
			require.Equal(t, *value, 7)
			require.Equal(t, isExisted, true)
			wg.Done()
		}()
		wg.Wait()

		wg.Add(1)
		go func() {
			cache.Put(4, 9) // Delete Least Recently Used
			wg.Done()
		}()
		wg.Wait()

		var b bytes.Buffer
		cache.PrintStatics(&b, cache.S, length) // statistics
		want := "Hit rate: 100.00\nEntries were written to the cache and have never been read: 2\nAverage number of times that things currently in the cache is read: 0.67\nTotal reads and writes have been performed in the cache including evicted: 6\n"
		require.Equal(t, b.String(), want)
		require.Equal(t, cache.S.EntriesNeverRead, 2)
		require.Equal(t, float32(cache.S.TotalReadExisted)/float32(len(cache.M)), float32(0.6666667))
		require.Equal(t, cache.S.ReadCount+cache.S.WritesCount, 6)

		wg.Add(1)
		go func() {
			_, isExisted := cache.Get(3)
			require.Equal(t, isExisted, false)
			wg.Done()
		}()
		wg.Wait()

		wg.Add(1)
		go func() {
			cache.Put(1, 11) // Update a value in the cache
			wg.Done()
		}()
		wg.Wait()
		
		value, isExisted := cache.Get(1)
		require.Equal(t, *value, 11)
		require.Equal(t, isExisted, true)
	})
}
