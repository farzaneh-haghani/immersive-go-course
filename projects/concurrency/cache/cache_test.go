package cache

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("Test a cache", func(t *testing.T) {
		length := 2
		cache := NewCache[string, string](length) // Create a cache
		require.Equal(t, cache.Size, length)

		cache.Put("key1", "value1") // Put in the cache

		got, isExisted := cache.Get("key1") // Get from the cache
		require.Equal(t, *got, "value1")
		require.Equal(t, isExisted, true)

		cache.Put("key2", "value2")

		got, isExisted = cache.Get("key2")
		require.Equal(t, *got, "value2")
		require.Equal(t, isExisted, true)

		cache.Put("key3", "value3") // Delete Least Recently Used

		var b bytes.Buffer
		cache.PrintStatics(&b, cache.S, length) // statistics
		want := "Hit rate: 100.00\nEntries were written to the cache and have never been read: 1\nAverage number of times that things currently in the cache is read: 0.50\nTotal reads and writes have been performed in the cache including evicted: 5\n"
		require.Equal(t, b.String(), want)
		require.Equal(t, cache.S.EntriesNeverRead, 1)
		require.Equal(t, float32(cache.S.TotalReadExisted)/float32(len(cache.M)), float32(0.5))
		require.Equal(t, cache.S.ReadCount+cache.S.WritesCount, 5)

		_, isExisted = cache.Get("key1")
		require.Equal(t, isExisted, false)

		got, isExisted = cache.Get("key2")
		require.Equal(t, *got, "value2")
		require.Equal(t, isExisted, true)

		got, isExisted = cache.Get("key3")
		require.Equal(t, *got, "value3")
		require.Equal(t, isExisted, true)

		cache.Put("key2", "value22") // Update a value in the cache
		got, isExisted = cache.Get("key2")
		require.Equal(t, *got, "value22")
		require.Equal(t, isExisted, true)
	})
}
