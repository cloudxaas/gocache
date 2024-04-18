package lrubytes

import (
    "time"
    "github.com/cloudxaas/gocx"
)

// CacheTTL struct includes the TTL for each entry
type CacheTTL struct {
    Cache
    ttlMap map[string]time.Time // Map to keep track of the TTL for each key
}

type entryTTL struct {
    entry
    expireAt time.Time // Time at which the entry expires
}

// NewLRUCacheTTL creates a new LRU Cache with TTL functionality
func NewLRUCacheTTL(maxMemory int64, evictBatchSize int, defaultTTL time.Duration) *CacheTTL {
    return &CacheTTL{
        Cache: Cache{
            maxMemory:      maxMemory,
            evictBatchSize: evictBatchSize,
            entries:        make([]entry, 0),
            indexMap:       make(map[string]int),
            head:           -1,
            tail:           -1,
        },
        ttlMap: make(map[string]time.Time),
    }
}

func (c *CacheTTL) Get(key []byte) ([]byte, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    keyStr := cx.B2s(key)
    if idx, ok := c.indexMap[keyStr]; ok {
        if time.Now().After(c.entries[idx].(entryTTL).expireAt) {
            c.Delete(key) // Expire the entry
            return nil, false
        }
        if idx != c.head {
            c.moveToFront(idx)
        }
        return c.entries[idx].value, true
    }
    return nil, false
}

func (c *CacheTTL) Put(key, value []byte, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    keyStr := cx.B2s(key)
    memSize := c.estimateMemory(key, value)
    expireAt := time.Now().Add(ttl)

    if idx, ok := c.indexMap[keyStr]; ok {
        oldMemSize := c.estimateMemory(c.entries[idx].key, c.entries[idx].value)
        c.adjustMemory(memSize - oldMemSize)
        c.entries[idx] = entryTTL{entry: entry{key: key, value: value}, expireAt: expireAt}
        c.moveToFront(idx)
        c.ttlMap[keyStr] = expireAt
        return
    }

    if c.currentMemory + memSize > c.maxMemory {
        for c.currentMemory + memSize > c.maxMemory {
            c.evict()
        }
    }

    idx := len(c.entries)
    c.entries = append(c.entries, entryTTL{entry: entry{key: key, value: value}, expireAt: expireAt})
    c.indexMap[keyStr] = idx
    c.adjustMemory(memSize)
    c.moveToFront(idx)
    c.ttlMap[keyStr] = expireAt
}

func (c *CacheTTL) evict() {
    for i := 0; i < c.evictBatchSize && c.tail != -1; i++ {
        tailKey := string(c.entries[c.tail].key)
        if _, ok := c.ttlMap[tailKey]; ok {
            delete(c.ttlMap, tailKey) // Also delete from ttlMap
        }
        c.detach(c.tail)
    }
}

func (c *CacheTTL) Delete(key []byte) {
    c.mu.Lock()
    defer c.mu.Unlock()
    keyStr := cx.B2s(key)
    if idx, ok := c.indexMap[keyStr]; ok {
        // Get the memory size of the entry to be deleted for proper memory adjustment
        memSize := c.estimateMemory(c.entries[idx].key, c.entries[idx].value)
        c.adjustMemory(-memSize) // Adjust the current memory usage
        c.detach(idx)            // Detach the entry from the linked list
        delete(c.indexMap, keyStr) // Remove the entry from the index map
        delete(c.ttlMap, keyStr)   // Remove the TTL entry
    }
}
