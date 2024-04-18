package lrubytes

import (
    "sync"
    cx "github.com/cloudxaas/gocx"
)

type Cache struct {
    maxMemory   int64
    currentMemory int64
    entries    []entry
    indexMap   map[string]int
    head, tail int
    mu         sync.Mutex
}

type entry struct {
    key, value []byte
    prev, next int
}

func NewLRUCache(maxMemory int64) *Cache {
    return &Cache{
        maxMemory: maxMemory,
        entries:   make([]entry, 0),  // Dynamic resizing based on memory usage
        indexMap:  make(map[string]int),
        head:      -1,
        tail:      -1,
    }
}

func (c *Cache) estimateMemory(key, value []byte) int64 {
    return int64(len(key) + len(value))
}

func (c *Cache) adjustMemory(delta int64) {
    c.currentMemory += delta
}

func (c *Cache) Get(key []byte) ([]byte, bool) {
    c.mu.Lock()
    keyStr := cx.B2s(key)
    if idx, ok := c.indexMap[keyStr]; ok {
        if idx != c.head {
            c.moveToFront(idx)
        }
        c.mu.Unlock()
        return c.entries[idx].value, true
    }
    c.mu.Unlock()
    return nil, false
}

func (c *Cache) Put(key, value []byte) {
    c.mu.Lock()
    keyStr := cx.B2s(key)  
    memSize := c.estimateMemory(key, value)

    if idx, ok := c.indexMap[keyStr]; ok {
        oldMemSize := c.estimateMemory(c.entries[idx].key, c.entries[idx].value)
        c.adjustMemory(memSize - oldMemSize)
        c.entries[idx].value = value
        c.moveToFront(idx)
        c.mu.Unlock()
        return
    }

    if c.currentMemory+memSize > c.maxMemory {
        // Evict least recently used items until there is enough space
        for c.currentMemory+memSize > c.maxMemory {
            c.evict()
        }
    }

    // Add new entry
    c.entries = append(c.entries, entry{key: key, value: value})
    idx := len(c.entries) - 1
    c.indexMap[keyStr] = idx
    c.adjustMemory(memSize)
    c.moveToFront(idx)
    c.mu.Unlock()
}

func (c *Cache) Delete(key []byte) {
    c.mu.Lock()
    keyStr := cx.B2s(key)
    if idx, ok := c.indexMap[keyStr]; ok {
        memSize := c.estimateMemory(c.entries[idx].key, c.entries[idx].value)
        c.adjustMemory(-memSize)
        c.detach(idx)
        delete(c.indexMap, keyStr)
    }
    c.mu.Unlock()
}

func (c *Cache) evict() {
    if c.tail != -1 {
        oldKeyStr := string(c.entries[c.tail].key)
        memSize := c.estimateMemory(c.entries[c.tail].key, c.entries[c.tail].value)
        c.adjustMemory(-memSize)
        c.detach(c.tail)
        delete(c.indexMap, oldKeyStr)
    }
}

func (c *Cache) detach(idx int) {
    if c.entries[idx].prev != -1 {
        c.entries[c.entries[idx].prev].next = c.entries[idx].next
    }
    if c.entries[idx].next != -1 {
        c.entries[c.entries[idx].next].prev = c.entries[idx].prev
    }
    if idx == c.head {
        c.head = c.entries[idx].next
    }
    if idx == c.tail {
        c.tail = c.entries[idx].prev
    }
}

func (c *Cache) moveToFront(idx int) {
    if idx == c.head {
        return
    }
    c.detach(idx)
    // Attach to front
    c.entries[idx].prev = -1
    c.entries[idx].next = c.head
    if c.head != -1 {
        c.entries[c.head].prev = idx
    }
    c.head = idx
}
