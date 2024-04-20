package main

import (
	"fmt"
	"github.com/cloudxaas/gocache/lrux/bytes"
)

// FNV-1a hash function for byte slices
func FNV1aHash(key []byte) uint32 {
	const (
		offset32 uint32 = 2166136261
		prime32  uint32 = 16777619
	)
	var hash uint32 = offset32
	for _, c := range key {
		hash ^= uint32(c)
		hash *= prime32
	}
	return hash
}

// Test integrity of the cache operations
func testCacheIntegrity(cache *lruxbytes.Cache) {
	fmt.Println("Testing single cache integrity...")

	// Set operations
	keys := [][]byte{[]byte("apple"), []byte("banana"), []byte("cherry")}
	values := [][]byte{[]byte("red"), []byte("yellow"), []byte("dark red")}

	for i, key := range keys {
		cache.Set(key, values[i])
	}

	// Get operations
	for i, key := range keys {
		value, ok := cache.Get(key)
		if ok && string(value) == string(values[i]) {
			fmt.Printf("Get operation successful for key %s, value %s\n", key, value)
		} else {
			fmt.Printf("Get operation failed for key %s\n", key)
		}
	}

	// Delete operation
	cache.Del(keys[0])
	if value, ok := cache.Get(keys[0]); !ok {
		fmt.Printf("Delete operation successful for key %s\n", keys[0])
	} else {
		fmt.Printf("Delete operation failed for key %s, returned value %s\n", keys[0], value)
	}
}

// Test integrity of the sharded cache operations
func testShardedCacheIntegrity(cache *lruxbytes.ShardedCache) {
	fmt.Println("Testing sharded cache integrity...")

	// Set operations
	keys := [][]byte{[]byte("mango"), []byte("kiwi"), []byte("grape")}
	values := [][]byte{[]byte("green"), []byte("brown"), []byte("purple")}

	for i, key := range keys {
		cache.Set(key, values[i])
	}

	// Get operations
	for i, key := range keys {
		value, ok := cache.Get(key)
		if ok && string(value) == string(values[i]) {
			fmt.Printf("Get operation successful for key %s, value %s\n", key, value)
		} else {
			fmt.Printf("Get operation failed for key %s\n", key)
		}
	}

	// Delete operation
	cache.Del(keys[1])
	if value, ok := cache.Get(keys[1]); !ok {
		fmt.Printf("Delete operation successful for key %s\n", keys[1])
	} else {
		fmt.Printf("Delete operation failed for key %s, returned value %s\n", keys[1], value)
	}
}

func main() {
	// Single cache example
	singleCache := lruxbytes.NewLRUCache(1024*100, 1, FNV1aHash) // 10 MB max, 1 item eviction batch
	testCacheIntegrity(singleCache)

	// Sharded cache example
	shardedCache := lruxbytes.NewShardedCache(16, 1024*100, 1, FNV1aHash) // 16 shards, 10 MB total, 1 item eviction batch
	testShardedCacheIntegrity(shardedCache)
}

