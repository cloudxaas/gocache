package cxcachelru

import (
    "sync"
)

// Sizer interface that requires a Size method, returning the size of the object in bytes.
type Sizer interface {
    Size() int64
}

// Define a type constraint that combines comparable and Sizer for keys.
type Key interface {
    comparable
    Sizer
}

// Cache struct definition using generic types K and V.
type Cache[K Key, V Sizer] struct {
    maxMemory       int64
    currentMemory   int64
    evictBatchSize  int
    entries         []entry[K, V]
    freeEntries     []int // stack of free entries indices
    indexMap        map[K]int
    head, tail      int
    mu              sync.Mutex
}

type entry[K Key, V Sizer] struct {
    key   K
    value V
    prev  int
    next  int
}

// NewLRUCache creates a new cache with specified maximum memory and eviction batch size.
func NewLRUCache[K Key, V Sizer](maxMemory int64, evictBatchSize int) *Cache[K, V] {
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

func (c *Cache[K, V]) estimateMemory(key K, value V) int64 {
    return key.Size() + value.Size()
}

func (c *Cache[K, V]) adjustMemory(delta int64) {
    c.currentMemory += delta
}

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

func (c *Cache[K, V]) Put(key K, value V) {
    c.mu.Lock()
    defer c.mu.Unlock()

    memSize := c.estimateMemory(key, value)

    if idx, ok := c.indexMap[key]; ok {
        oldMemSize := c.estimateMemory(c.entries[idx].key, c.entries[idx].value)
        c.adjustMemory(memSize - oldMemSize)
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
    } else {
        idx = len(c.entries)
        c.entries = append(c.entries, entry[K, V]{})
    }

    c.entries[idx] = entry[K, V]{key: key, value: value, prev: -1, next: -1}
    c.indexMap[key] = idx
    c.adjustMemory(memSize)
    c.moveToFront(idx)
}
