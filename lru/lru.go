package cxcachelru

import (
    "sync"
)

// Sizer interface requires a method Size that returns the size of the object in bytes.
type Sizer interface {
    Size() int64
}

// Cache struct definition with generics.
// K and V are the types for the key and value used in the cache. They must satisfy the Sizer interface.
type Cache[K Sizer, V Sizer] struct {
    maxMemory       int64
    currentMemory   int64
    evictBatchSize  int
    entries         []entry[K, V]
    freeEntries     []int // Stack of indices of free entries
    indexMap        map[K]int
    head, tail      int
    mu              sync.Mutex
}

// Entry holds a key, a value, and pointers to other entries in the LRU cache.
type entry[K Sizer, V Sizer] struct {
    key   K
    value V
    prev  int
    next  int
}

// NewLRUCache creates a new LRU Cache with specified max memory and eviction batch size.
func NewLRUCache[K Sizer, V Sizer](maxMemory int64, evictBatchSize int) *Cache[K, V] {
    return &Cache[K, V]{
        maxMemory:      maxMemory,
        evictBatchSize: evictBatchSize,
        entries:        make([]entry[K, V], 0),
        indexMap:       make(map[K]int),
        head:           -1,
        tail:           -1,
        freeEntries:    make([]int, 0),
    }
}

// estimateMemory calculates the total memory usage for the key-value pair.
func (c *Cache[K, V]) estimateMemory(key K, value V) int64 {
    return key.Size() + value.Size()
}

// Get retrieves the value for a key from the cache.
func (c *Cache[K, V]) Get(key K) (V, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if idx, ok := c.indexMap[key]; ok {
        if idx != c.head {
            c.moveToFront(idx)
        }
        return c.entries[idx].value, true
    }
    var zero V
    return zero, false
}

// Put adds a key-value pair to the cache, managing memory and evicting as necessary.
func (c *Cache[K, V]) Put(key K, value V) {
    c.mu.Lock()
    defer c.mu.Unlock()

    memSize := c.estimateMemory(key, value)

    if idx, ok := c.indexMap[key]; ok {
        c.adjustMemory(memSize) // Assuming new and old values have same estimated memory
        c.entries[idx].value = value
        c.moveToFront(idx)
        return
    }

    if c.currentMemory + memSize > c.maxMemory {
        c.evict()
    }

    var idx int
    if len(c.freeEntries) > 0 {
        idx = c.freeEntries[len(c.freeEntries)-1]
        c.freeEntries = c.freeEntries[:len(c.freeEntries)-1]
        c.entries[idx] = entry[K, V]{key: key, value: value}
    } else {
        c.entries = append(c.entries, entry[K, V]{key: key, value: value})
        idx = len(c.entries) - 1
    }

    c.indexMap[key] = idx
    c.adjustMemory(memSize)
    c.moveToFront(idx)
}

// moveToFront updates the cache to move a given index to the front (most recently used).
func (c *Cache[K, V]) moveToFront(idx int) {
    if idx == c.head {
        return
    }
    c.detach(idx)

    if c.head != -1 {
        c.entries[c.head].prev = idx
    }
    c.entries[idx].next = c.head
    c.entries[idx].prev = -1
    c.head = idx

    if c.tail == -1 {
        c.tail = idx
    }

    if c.tail == idx {
        c.tail = c.entries[idx].prev
    }
}

// Delete removes a key from the cache.
func (c *Cache[K, V]) Delete(key K) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if idx, ok := c.indexMap[key]; ok {
        c.adjustMemory(-c.estimateMemory(c.entries[idx].key, c.entries[idx].value))
        c.detach(idx)
        c.freeEntries = append(c.freeEntries, idx)
        delete(c.indexMap, key)
    }
}

// detach removes an entry from the linked list part of the cache.
func (c *Cache[K, V]) detach(idx int) {
    if c.entries[idx].prev != -1 {
        c.entries[c.entries[idx].prev].next = c.entries[idx].next
    } else {
        c.head = c.entries[idx].next
    }

    if c.entries[idx].next != -1 {
        c.entries[c.entries[idx].next].prev = c.entries[idx].prev
    } else {
        c.tail = c.entries[idx].prev
    }

    if c.head == -1 {
        c.tail = -1
    }
}

// evict removes the least recently used items based on the eviction batch size.
func (c *Cache[K, V]) evict() {
    for i := 0; i < c.evictBatchSize && c.tail != -1; i++ {
        idx := c.tail
        c.adjustMemory(-c.estimateMemory(c.entries[idx].key, c.entries[idx].value))
        c.detach(idx)
        c.freeEntries = append(c.freeEntries, idx)
    }
}

// adjustMemory modifies the current memory tracking.
func (c *Cache[K, V]) adjustMemory(delta int64) {
    c.currentMemory += delta
}
