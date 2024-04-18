package lrubytes

import (
    "testing"
)

func BenchmarkPut(b *testing.B) {
    cache := NewLRUCache(1024 * 1024 * 10) // 10 MB max memory
    keys := make([][]byte, 1000)
    values := make([][]byte, 1000)
    for i := 0; i < 1000; i++ {
        keys[i] = []byte{byte(i)}
        values[i] = make([]byte, 1024) // 1 KB values
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Put(keys[i%1000], values[i%1000])
    }
}

func BenchmarkGet(b *testing.B) {
    cache := NewLRUCache(1024 * 1024 * 10) // 10 MB max memory
    for i := 0; i < 1000; i++ {
        cache.Put([]byte{byte(i)}, make([]byte, 1024)) // 1 KB values
    }
    keys := make([][]byte, b.N)
    for i := 0; i < b.N; i++ {
        keys[i] = []byte{byte(i % 1000)}
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = cache.Get(keys[i])
    }
}

func BenchmarkDelete(b *testing.B) {
    cache := NewLRUCache(1024 * 1024 * 10) // 10 MB max memory
    keys := make([][]byte, 1000)
    for i := 0; i < 1000; i++ {
        keys[i] = []byte{byte(i)}
        cache.Put(keys[i], make([]byte, 1024)) // 1 KB values
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Delete(keys[i%1000])
    }
}


// Benchmarking the sharded cache version with concurrent access
func BenchmarkShardedPutParallel(b *testing.B) {
    shardCount := uint8(16)  // Example shard count
    totalMemory := int64(10 * 1024 * 1024) // 10 MB total memory
    cache := NewShardedCache(shardCount, totalMemory)

    keys := make([][]byte, 1000)
    values := make([][]byte, 1000)
    for i := 0; i < 1000; i++ {
        keys[i] = []byte{byte(i)}
        values[i] = make([]byte, 1024) // 1 KB values
    }

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for i := 0; pb.Next(); i++ {
            cache.Put(keys[i%1000], values[i%1000])
        }
    })
}

func BenchmarkShardedGetParallel(b *testing.B) {
    shardCount := uint8(16)  // Example shard count
    totalMemory := int64(10 * 1024 * 1024) // 10 MB total memory
    cache := NewShardedCache(shardCount, totalMemory)

    // Pre-fill the cache to ensure there's something to get
    for i := 0; i < 1000; i++ {
        cache.Put([]byte{byte(i)}, make([]byte, 1024)) // 1 KB values
    }

    keys := make([][]byte, 1000)
    for i := 0; i < 1000; i++ {
        keys[i] = []byte{byte(i)}
    }

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for i := 0; pb.Next(); i++ {
            _, _ = cache.Get(keys[i%1000])
        }
    })
}

func BenchmarkShardedDeleteParallel(b *testing.B) {
    shardCount := uint8(16)  // Example shard count
    totalMemory := int64(10 * 1024 * 1024) // 10 MB total memory
    cache := NewShardedCache(shardCount, totalMemory)

    // Pre-fill the cache to ensure there's something to delete
    keys := make([][]byte, 1000)
    for i := 0; i < 1000; i++ {
        keys[i] = []byte{byte(i)}
        cache.Put(keys[i], make([]byte, 1024)) // 1 KB values
    }

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for i := 0; pb.Next(); i++ {
            cache.Delete(keys[i%1000])
        }
    })
}
