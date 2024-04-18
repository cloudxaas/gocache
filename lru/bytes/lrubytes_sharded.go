package lrubytes

import (
    "github.com/zeebo/xxh3"
)

// ShardedCache struct containing multiple Cache shards
type ShardedCache struct {
    shards     []*Cache
    shardCount uint8
}

// NewShardedCache creates a new ShardedCache with the specified number of shards and total memory limit
func NewShardedCache(shardCount uint8, totalMemory int64) *ShardedCache {
    maxMemoryPerShard := totalMemory / int64(shardCount) // Calculate memory per shard
    shards := make([]*Cache, shardCount)
    for i := uint8(0); i < shardCount; i++ {
        shards[i] = NewLRUCache(maxMemoryPerShard)
    }
    return &ShardedCache{
        shards:     shards,
        shardCount: shardCount,
    }
}

// getShard computes the hash of the key to determine which shard to use
func (sc *ShardedCache) getShard(key []byte) *Cache {
    hash := xxh3.Hash(key)
    return sc.shards[uint8(hash)%sc.shardCount]
}

// Get retrieves a value from the appropriate shard
func (sc *ShardedCache) Get(key []byte) ([]byte, bool) {
    shard := sc.getShard(key)
    return shard.Get(key)
}

// Put adds a key-value pair to the appropriate shard
func (sc *ShardedCache) Put(key, value []byte) {
    shard := sc.getShard(key)
    shard.Put(key, value)
}

// Delete removes a key from the appropriate shard
func (sc *ShardedCache) Delete(key []byte) {
    shard := sc.getShard(key)
    shard.Delete(key)
}
