package cache

import (
	"sync"

	"github.com/AnnDutova/otus_go_hw/hw04_lru_cache/list"
)

type Cache interface {
	Set(key list.Key, value interface{}) bool
	Get(key list.Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    list.List
	items    map[list.Key]*list.Item
	m        sync.Mutex
}

func (c *lruCache) Set(key list.Key, value interface{}) bool {
	c.m.Lock()
	defer c.m.Unlock()
	_, ok := c.items[key]

	if ok {
		newElem := &list.Item{Value: value}
		c.items[key] = newElem
		c.queue.MoveToFront(newElem)
		return ok
	}

	if c.queue.Len() == c.capacity {
		delete(c.items, c.queue.Back().Key)
		c.queue.Remove(c.queue.Back())
	}

	newElem := c.queue.PushFront(value)
	newElem.Key = key
	c.items[key] = newElem

	return ok
}

func (c *lruCache) Get(key list.Key) (interface{}, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	if el, ok := c.items[key]; ok {
		c.queue.MoveToFront(el)
		return el.Value, ok
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.m.Lock()
	defer c.m.Unlock()

	c.queue = list.NewList()
	c.items = map[list.Key]*list.Item{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    list.NewList(),
		items:    make(map[list.Key]*list.Item, capacity),
		m:        sync.Mutex{},
	}
}
