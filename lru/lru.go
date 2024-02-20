package cxcachelru

import (
	"container/list"
	"sync"
)

// Entry is a key-value pair stored in the cache.
type Entry struct {
	key   interface{}
	value interface{}
}

// LRUCache is a zero allocation, zero GC, no heap escape LRU cache.
type LRUCache struct {
	capacity int
	cache    map[interface{}]*list.Element
	list     *list.List
	mutex    sync.Mutex
}

// NewLRUCache creates a new LRUCache with the given capacity.
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[interface{}]*list.Element, capacity),
		list:     list.New(),
	}
}

// Get retrieves a value from the cache.
func (c *LRUCache) Get(key interface{}) (value interface{}, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, hit := c.cache[key]; hit {
		c.list.MoveToFront(elem)
		return elem.Value.(*Entry).value, true
	}
	return nil, false
}

// Put adds a key-value pair to the cache.
func (c *LRUCache) Put(key, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, hit := c.cache[key]; hit {
		c.list.MoveToFront(elem)
		elem.Value.(*Entry).value = value
		return
	}

	if c.list.Len() == c.capacity {
		// Evict the least recently used (LRU) element.
		lastElem := c.list.Back()
		if lastElem != nil {
			c.list.Remove(lastElem)
			delete(c.cache, lastElem.Value.(*Entry).key)
		}
	}

	// Use existing entry object if available to avoid allocations.
	entry := &Entry{key, value}
	elem := c.list.PushFront(entry)
	c.cache[key] = elem
}

// Size returns the number of items in the cache.
func (c *LRUCache) Size() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.list.Len()
}
