package cache

import (
	"sync"

	"github.com/AnnDutova/otus_go_hw/hw04_lru_cache/list"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    list.List
	items    map[Key]*list.Item
	m        sync.Mutex
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.m.Lock()
	defer c.m.Unlock()
	_, ok := c.items[key]

	if ok {
		newElem := &list.Item{Value: value}
		c.items[key] = newElem
		c.queue.MoveToFront(newElem)
		return ok
	}

	newElem := &list.Item{Value: value}
	c.items[key] = newElem
	if c.queue.Len()+1 > c.capacity {
		tail := c.queue.Back()
		c.queue.Remove(tail)

		var deleteKey Key
		for k, v := range c.items {
			if v.Value == tail.Value {
				deleteKey = k
				break
			}
		}
		delete(c.items, deleteKey)

		c.queue.PushFront(value)
	} else {
		c.queue.PushFront(value)
	}
	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
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
	c.items = map[Key]*list.Item{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    list.NewList(),
		items:    make(map[Key]*list.Item, capacity),
		m:        sync.Mutex{},
	}
}
