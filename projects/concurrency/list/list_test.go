package list

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddNode(t *testing.T) {

	t.Run("Test a link list", func(t *testing.T) {
		l := NewList[int, string]()
		require.Equal(t, l.first, l.last)

		newNode1 := l.AddNode(1, "test1")
		require.Equal(t, newNode1.key, 1)
		require.Equal(t, newNode1.Value, "test1")
		require.Equal(t, l.first, newNode1)
		require.Equal(t, l.first, l.last)
		require.Equal(t, newNode1.next, newNode1.prev)

		newNode2 := l.AddNode(2, "test2")
		require.Equal(t, newNode2.key, 2)
		require.Equal(t, newNode2.Value, "test2")
		require.Equal(t, l.first, newNode1)
		require.Equal(t, l.last, newNode2)
		require.Equal(t, newNode1.next, newNode2)
		require.Equal(t, newNode2.prev, newNode1)
		require.Equal(t, newNode2.next, newNode1.prev)

		l.MoveNodeToLast(newNode1)
		require.Equal(t, l.last, newNode1)
		require.Equal(t, l.first, newNode2)
		require.Equal(t, newNode2.next, newNode1)
		require.Equal(t, newNode1.prev, newNode2)
		require.Equal(t, newNode2.prev, newNode1.next)

		deleted, _ := l.DeleteFirstNode()
		require.Equal(t, deleted, 2)
		require.NotEqual(t, l.first, newNode2)
		require.Equal(t, l.first, newNode1)
		require.Equal(t, l.first, l.last)
		require.Equal(t, newNode1.next, newNode1.prev)

		deleted2, _ := l.DeleteFirstNode()
		require.Equal(t, deleted2, 1)
		require.NotEqual(t, l.first, newNode1)
		require.Equal(t, l.first, l.last)
	})
}
