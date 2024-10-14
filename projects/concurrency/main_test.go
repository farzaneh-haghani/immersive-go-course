package main

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("Test as a non-thread safe cache", func(t *testing.T) {

		cache1 := NewCache[string, string](2) // Create a cache
		require.Equal(t, cache1.size, 2)

		cache1.Put("key1", "value1") // Put in the cache

		got, isExisted := cache1.Get("key1") // Get from the cache
		require.Equal(t, *got, "value1")
		require.Equal(t, isExisted, true)

		cache1.Put("key2", "value2")

		got, isExisted = cache1.Get("key2")
		require.Equal(t, *got, "value2")
		require.Equal(t, isExisted, true)

		cache1.Put("key3", "value3") // Delete Least Recently Used

		hit := hitRate(cache1.s.readCount, cache1.s.unreadCount) // statistics
		require.Equal(t, hit, float32(100))
		require.Equal(t, cache1.s.entriesNeverRead, 1)
		require.Equal(t, float32(cache1.s.totalReadExisted)/float32(len(cache1.m)), float32(0.5))
		require.Equal(t, cache1.s.readCount+cache1.s.writesCount, 5)

		_, isExisted = cache1.Get("key1")
		require.Equal(t, isExisted, false)

		got, isExisted = cache1.Get("key2")
		require.Equal(t, *got, "value2")
		require.Equal(t, isExisted, true)

		got, isExisted = cache1.Get("key3")
		require.Equal(t, *got, "value3")
		require.Equal(t, isExisted, true)

		cache1.Put("key2", "value22") // Update a value in the cache
		got, isExisted = cache1.Get("key2")
		require.Equal(t, *got, "value22")
		require.Equal(t, isExisted, true)
	})

	t.Run("Test as a thread safe cache", func(t *testing.T) {
		var wg sync.WaitGroup

		cache2 := NewCache[int, int](3) // Create a cache

		for i := 1; i < 4; i++ {
			wg.Add(1)
			go func() {
				cache2.Put(i, i+5) // Put in the cache by 3 go routines
				wg.Done()
			}()
		}

		wg.Wait()
		wg.Add(2)
		go func() {
			value, isExisted := cache2.Get(1) // Get from the cache
			require.Equal(t, *value, 6)
			require.Equal(t, isExisted, true)
			wg.Done()
		}()

		go func() {
			value, isExisted := cache2.Get(2)
			require.Equal(t, *value, 7)
			require.Equal(t, isExisted, true)
			wg.Done()
		}()
		wg.Wait()

		wg.Add(1)
		go func() {
			cache2.Put(4, 9) // Delete Least Recently Used
			wg.Done()
		}()
		wg.Wait()

		hit := hitRate(cache2.s.readCount, cache2.s.unreadCount) // statistics
		require.Equal(t, hit, float32(100))
		require.Equal(t, cache2.s.entriesNeverRead, 2)
		require.Equal(t, float32(cache2.s.totalReadExisted)/float32(len(cache2.m)), float32(0.6666667))
		require.Equal(t, cache2.s.readCount+cache2.s.writesCount, 6)

		wg.Add(1)
		go func() {
			_, isExisted := cache2.Get(3)
			require.Equal(t, isExisted, false)
			wg.Done()
		}()
		wg.Wait()

		wg.Add(1)
		go func() {
			cache2.Put(1, 11) // Update a value in the cache
			wg.Done()
		}()
		wg.Wait()
		
		value, isExisted := cache2.Get(1)
		require.Equal(t, *value, 11)
		require.Equal(t, isExisted, true)
	})
}
