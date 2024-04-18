# Fastest Zero Allocation LRU Cache for Go key value in []byte 

Welcome to the repository for the fastest LRU cache implementations available for Go. This LRU cache is uniquely designed to limit the memory usage directly, rather than just the number of entries. This makes it ideal for applications where the memory footprint is critical, such as in embedded systems or high-performance computing environments where resources are tightly managed.

## Features

- **Memory-Size Limited**: Unlike other LRU caches that limit the number of entries, this cache controls the total memory used, allowing for better resource management in memory-constrained environments.
- **High Performance**: Designed with performance in mind, benchmarks demonstrate extremely low latency and zero allocations during operations, ensuring minimal impact on application throughput.
- **Concurrency Safe**: Implements synchronization to manage concurrent access, making it suitable for high-concurrency scenarios.

## Benchmarks

The cache has been rigorously benchmarked on a system with the following specifications:
- **OS**: Linux
- **Architecture**: AMD64
- **CPU**: AMD Ryzen 5 7640HS w/ Radeon 760M Graphics

Benchmark results:
```
go test -bench=. -benchmem
goos: linux
goarch: amd64
pkg: github.com/cloudxaas/gocache/lru/bytes
cpu: AMD Ryzen 5 7640HS w/ Radeon 760M Graphics     
BenchmarkPut-12                      	74357374	        15.97 ns/op	       0 B/op	       0 allocs/op
BenchmarkGet-12                      	70288532	        14.68 ns/op	       0 B/op	       0 allocs/op
BenchmarkDelete-12                   	235118634	         5.088 ns/op	       0 B/op	       0 allocs/op
BenchmarkPutParallel-12              	21257577	        52.22 ns/op	       0 B/op	       0 allocs/op
BenchmarkGetParallel-12              	24624946	        48.31 ns/op	       0 B/op	       0 allocs/op
BenchmarkDeleteParallel-12           	36608547	        29.73 ns/op	       0 B/op	       0 allocs/op
BenchmarkShardedPutParallel-12       	48741004	        24.90 ns/op	       0 B/op	       0 allocs/op
BenchmarkShardedGetParallel-12       	50258973	        23.46 ns/op	       0 B/op	       0 allocs/op
BenchmarkShardedDeleteParallel-12    	166783006	         6.580 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/cloudxaas/gocache/lru/bytes	14.711s
```

These benchmarks illustrate the efficiency and speed of the cache, which is designed to operate with zero memory allocations during runtime operations, contributing to its high performance.

## Usage

To use this cache, include it in your Go project and create a cache instance specifying the maximum memory it should use:

```go
package main

import (
    cxlrubytes "github.com/cloudxaas/gocache/lru/bytes"
    "fmt"
)

func main() {
    // Initialize a new LRU cache with a max memory limit of 10 MB
    cache := cxlrubytes.NewLRUCache(10 * 1024 * 1024)

    // Example of adding a value to the cache
    cache.Put([]byte("key1"), []byte("value1"))

    // Retrieve a value
    if value, found := cache.Get([]byte("key1")); found {
        fmt.Println("Retrieved:", string(value))
    }

    // Delete a value
    cache.Delete([]byte("key1"))
}
```


### Sharded version
```go
package main

import (
    cxlrubytes "github.com/cloudxaas/gocache/lru/bytes"
    "fmt"
)

func main() {
    // Initialize a new sharded LRU cache with a total memory limit of 10 MB across 16 shards
    shardCount := uint8(16)
    totalMemory := int64(10 * 1024 * 1024) // 10 MB total memory for the cache
    cache := cxlrubytes.NewShardedCache(shardCount, totalMemory)

    // Example of adding and retrieving values
    cache.Put([]byte("key1"), []byte("value1"))
    if value, found := cache.Get([]byte("key1")); found {
        fmt.Println("Retrieved:", string(value))
    }

    // Delete a value
    cache.Delete([]byte("key1"))
}
```
