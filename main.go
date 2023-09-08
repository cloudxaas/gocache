package main

import (
	"container/list"
	"fmt"
)

type LRUCache struct {
	capacity int
	arena    []CacheEntry
	cacheMap map[int]*list.Element
	lruList  *list.List
}

type CacheEntry struct {
	key   int
	value int
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		arena:    make([]CacheEntry, 0, capacity),
		cacheMap: make(map[int]*list.Element),
		lruList:  list.New(),
	}
}

func (c *LRUCache) Get(key int) (int, bool) {
	if element, ok := c.cacheMap[key]; ok {
		c.lruList.MoveToFront(element)
		entry := element.Value.(*CacheEntry)
		return entry.value, true
	}
	return -1, false
}

func (c *LRUCache) Put(key int, value int) {
	if element, ok := c.cacheMap[key]; ok {
		c.lruList.MoveToFront(element)
		element.Value.(*CacheEntry).value = value
		return
	}

	if c.lruList.Len() == c.capacity {
		// Remove the least recently used entry
		element := c.lruList.Back()
		removedEntry := c.lruList.Remove(element).(*CacheEntry)
		delete(c.cacheMap, removedEntry.key)
	} else {
		newEntry := &CacheEntry{key: key, value: value}
		c.arena = append(c.arena, *newEntry)
		newEntry = &c.arena[len(c.arena)-1]
		c.cacheMap[key] = c.lruList.PushFront(newEntry)
	}
}

func main() {
	cache := NewLRUCache(2)

	cache.Put(1, 1)
	cache.Put(2, 2)

	fmt.Println(cache.Get(1)) // Output: 1

	cache.Put(3, 3)

	fmt.Println(cache.Get(2)) // Output: -1

	cache.Put(4, 4)

	fmt.Println(cache.Get(1)) // Output: -1
	fmt.Println(cache.Get(3)) // Output: 3
	fmt.Println(cache.Get(4)) // Output: 4
}
