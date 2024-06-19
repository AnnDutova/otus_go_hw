package list_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/AnnDutova/otus_go_hw/hw04_lru_cache/list"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := list.NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("push front", func(t *testing.T) {
		l := list.NewList()

		l.PushFront(10) // [10]
		l.PushFront(20) // [20 10]

		require.Equal(t, 2, l.Len())
		require.NotNil(t, l.Front())
		require.Equal(t, l.Front().Value, 20)
		require.NotNil(t, l.Back())
		require.Equal(t, l.Back().Value, 10)

		l.PushFront(30) // [30 20 10]

		require.Equal(t, 3, l.Len())
		require.NotNil(t, l.Front())
		require.Equal(t, l.Front().Value, 30)
		require.NotNil(t, l.Back())
		require.Equal(t, l.Back().Value, 10)
	})

	t.Run("push back", func(t *testing.T) {
		l := list.NewList()

		l.PushBack(10) // [10]
		l.PushBack(20) // [10 20]

		require.Equal(t, 2, l.Len())
		require.NotNil(t, l.Front())
		require.Equal(t, l.Front().Value, 10)
		require.NotNil(t, l.Back())
		require.Equal(t, l.Back().Value, 20)

		l.PushBack(30) // [10 20 30]
		l.PrintList()

		require.Equal(t, 3, l.Len())
		require.NotNil(t, l.Front())
		require.Equal(t, l.Front().Value, 10)
		require.NotNil(t, l.Back())
		require.Equal(t, l.Back().Value, 30)
	})

	t.Run("remove", func(t *testing.T) {
		l := list.NewList()

		l.PushBack(10) // [10]
		l.PushBack(20) // [10 20]
		l.PushBack(30) // [10 20 30]
		l.PushBack(40) // [10 20 30 40]
		l.PrintList()
		require.Equal(t, 4, l.Len())
		require.Equal(t, l.Front().Value, 10)
		require.Equal(t, l.Back().Value, 40)

		l.Remove(l.Back()) // [10 20 30]
		l.PrintList()
		require.Equal(t, 3, l.Len())
		require.Equal(t, l.Front().Value, 10)
		require.Equal(t, l.Back().Value, 30)

		l.Remove(l.Front()) // [20 30]
		l.PrintList()
		require.Equal(t, 2, l.Len())
		require.Equal(t, l.Front().Value, 20)
		require.Equal(t, l.Back().Value, 30)

		value := l.Back() // 30
		l.PushBack(50)    // [20 30 50]
		l.PrintList()
		l.Remove(value)
		l.PrintList() // [20 50]
	})

	t.Run("complex", func(t *testing.T) {
		l := list.NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.PrintList()
		l.MoveToFront(l.Back()) // [70, 80, 60, 40, 10, 30, 50]
		l.PrintList()

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
