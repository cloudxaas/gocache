package lruxbytes

import (
	"fmt"
)

// ShardedCache struct containing multiple Cache shards
type ShardedCache struct {
	shards     []*Cache
	shardCount uint8
	hashFunc   ByteHashFunc
}

// NewShardedCache creates a new ShardedCache with the specified number of shards, total memory limit, eviction count, and a hash function
func NewShardedCache(shardCount uint8, totalMemory int64, evictionCount int, hashFunc ByteHashFunc) *ShardedCache {
	if shardCount == 0 || (shardCount&(shardCount-1)) != 0 {
		panic(fmt.Errorf("shardCount must be a non-zero power of 2, got %d", shardCount))
	}
	maxMemoryPerShard := totalMemory / int64(shardCount)
	shards := make([]*Cache, shardCount)
	for i := uint8(0); i < shardCount; i++ {
		shards[i] = NewLRUCache(maxMemoryPerShard, evictionCount, hashFunc)
	}
	return &ShardedCache{
		shards:     shards,
		shardCount: shardCount,
		hashFunc:   hashFunc,
	}
}

// getShard computes the hash of the key to determine which shard to use
func (sc *ShardedCache) getShard(key []byte) *Cache {
	hash := sc.hashFunc(key)
	return sc.shards[uint8(hash)&(sc.shardCount-1)]
}

// Get retrieves a value from the appropriate shard
func (sc *ShardedCache) Get(key []byte) ([]byte, bool) {
	shard := sc.getShard(key)
	return shard.Get(key)
}

// Set adds a key-value pair to the appropriate shard
func (sc *ShardedCache) Set(key, value []byte) {
	shard := sc.getShard(key)
	shard.Set(key, value)
}

// Del removes a key from the appropriate shard
func (sc *ShardedCache) Del(key []byte) {
	shard := sc.getShard(key)
	shard.Del(key)
}
