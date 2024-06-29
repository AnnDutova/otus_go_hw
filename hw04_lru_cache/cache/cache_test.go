package cache_test

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/AnnDutova/otus_go_hw/hw04_lru_cache/cache"
	"github.com/AnnDutova/otus_go_hw/hw04_lru_cache/list"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := cache.NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := cache.NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := cache.NewCache(3)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)
		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)
		wasInCache = c.Set("ccc", 300)
		require.False(t, wasInCache)
		wasInCache = c.Set("ddd", 400)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Equal(t, nil, val)
	})

	t.Run("oldest", func(t *testing.T) {
		c := cache.NewCache(3)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)
		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)
		wasInCache = c.Set("ccc", 300)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("ccc")
		require.True(t, ok)
		require.Equal(t, 300, val)

		wasInCache = c.Set("ddd", 400)
		require.False(t, wasInCache)

		val, ok = c.Get("bbb")
		require.False(t, ok)
		require.Equal(t, nil, val)
	})
}

func TestCacheMultithreading(_ *testing.T) {
	c := cache.NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(list.Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(list.Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
