package computingCache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComputing(t *testing.T) {
	myCache := NewComputingCache[int, string](3, func(int) string {
		return "test123"
	})

	myValue := myCache.Get(1)
	require.Equal(t, myValue, "test123")
	// myCache.cache.Put(2, "test222") // The reason that we didn't use embedding (myCache.Put) will work and being everything public isn't safe.
	myValue2 := myCache.Get(2)
	require.Equal(t, myValue2, "test123")
}
