<p align="center">

  <img src="https://github.com/cloudxaas/gocache/assets/104323920/5948a699-64c8-47b8-a5d6-5afedb6a3976" width="40%" height="auto">

    
 <h2 align="center">" Let's <u>Acce</u>lerate <u>LRU</u> " - Accelru</h2>
</p>

# Fastest Zero Allocation LRU Cache in Golang (for key, value pairs in []byte) - Accelru (cxlrubytes)

Welcome to the repository for the fastest LRU cache implementations available for Go. This LRU cache is uniquely designed to limit the memory usage directly, rather than by the number of entries. This makes it ideal for applications where the memory footprint is critical, such as in embedded systems or high-performance computing environments where resources are tightly managed.

(adjust the limit of eviction count to match usage scenario)
if you have roughly (you think) 1 mil items (capacity), you can set eviction count to 1024 or more. depending on usage patterns too. set the parameters accordingly.

## Features

- **Memory-Size Limited**: Unlike other LRU caches that limit the number of entries, this cache controls the total memory used, allowing for better resource management in memory-constrained environments.
- **High Performance**: Designed with performance in mind, benchmarks demonstrate extremely low latency and zero allocations during operations, ensuring minimal impact on application throughput.
- **Concurrency Safe**: Implements synchronization to manage concurrent access, making it suitable for high-concurrency scenarios.
- **1 Item or Batch evictions**: When cache is filled, eviction by batch or 1 by 1, you can set this value.

## Motivation
Most lru cache available online for golang are set by capacity count, which means you may OOM your program. With this lru, once you set the memory size limit, you do not need to worry about OOM or garbage collection issues with zero allocation. OOM = Out of memory.

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
BenchmarkPhusluLRUPut-12                       	65263994	        18.39 ns/op	       0 B/op	       0 allocs/op
BenchmarkPhusluLRUGet-12                       	63551188	        17.08 ns/op	       0 B/op	       0 allocs/op
BenchmarkPhusluLRUDelete-12                    	134234379	         8.743 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesPut-12                      	77965520	        15.65 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesGet-12                      	71972727	        14.94 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesDelete-12                   	295424088	         4.015 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesPutParallel-12              	23514391	        47.82 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesGetParallel-12              	24274286	        44.65 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesDeleteParallel-12           	40343742	        28.09 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesShardedPut-12               	50383508	        23.34 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesShardedGet-12               	51598164	        22.09 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesShardedDelete-12            	151449612	         7.935 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesShardedPutParallel-12       	45616530	        24.81 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesShardedGetParallel-12       	52305670	        23.15 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesShardedDeleteParallel-12    	220935474	         5.780 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/cloudxaas/gocache/lru/bytes	25.243s

```
reference : https://github.com/phuslu/lru


These benchmarks illustrate the efficiency and speed of the cache, which is designed to operate with zero memory allocations during runtime operations, contributing to its high performance.

## Usage

To use this cache, include it in your Go project and create a cache instance specifying the maximum memory it should use and the number of items to be evicted in a single go, which is faster than 1 by 1 eviction:


```go
package main

import (
    cxlrubytes "github.com/cloudxaas/gocache/lru/bytes"
    "fmt"
)

func main() {
    // Initialize a new LRU cache with a max memory limit of 10 MB, with an eviction count of 1024 at one go
    cache := cxlrubytes.NewLRUCache(10 * 1024 * 1024, 1024)

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
or

```go
package main

import (
    cx "github.com/cloudxaas/gocx"
    cxlrubytes "github.com/cloudxaas/gocache/lru/bytes"
    "fmt"
)

func main() {
    // Initialize a new LRU cache with a max memory limit of 10 MB, with an eviction count of 1024 at one go
    cache := cxlrubytes.NewLRUCache(10 * 1024 * 1024, 1024)

    // Example of adding a value to the cache
    cache.Put(cx.S2b("key1"), cx.S2b("value1"))

    // Retrieve a value
    if value, found := cache.Get(cx.S2b("key1")); found {
        fmt.Println("Retrieved:", cx.B2s(value))
    }

    // Delete a value
    cache.Delete(cx.S2b("key1"))
}
```


### Sharded version

Theoretically should work better in high concurrency environment with multiple goroutines.
Use this option when you have a lot of cpu cores.

```go

package main

import (
    cxlrubytes "github.com/cloudxaas/gocache/lru/bytes"
    "fmt"
)

func main() {
    // Initialize a new sharded LRU cache with a total memory limit of 10 MB across 16 shards
    shardCount := uint8(16)
    totalMemory := int64(10 * 1024 * 1024) // 10 MB total memory for the cache, with an eviction count of 1024 at one go
    cache := cxlrubytes.NewShardedCache(shardCount, totalMemory, 1024)

    // Example of adding and retrieving values
    cache.Put([]byte("key1"), []byte("value1"))
    if value, found := cache.Get([]byte("key1")); found {
        fmt.Println("Retrieved:", string(value))
    }

    // Delete a value
    cache.Delete([]byte("key1"))
}
```
or


```go
package main

import (
    cx "github.com/cloudxaas/gocx"
    cxlrubytes "github.com/cloudxaas/gocache/lru/bytes"
    "fmt"
)

func main() {
    // Initialize a new sharded LRU cache with a total memory limit of 10 MB across 16 shards
    shardCount := uint8(16)
    totalMemory := int64(10 * 1024 * 1024) // 10 MB total memory for the cache, with an eviction count of 1024 at one go
    cache := cxlrubytes.NewShardedCache(shardCount, totalMemory, 1024)

    // Example of adding and retrieving values
    cache.Put(cx.S2b("key1"), cx.S2b("value1"))
    if value, found := cache.Get(cx.S2b("key1")); found {
        fmt.Println("Retrieved:", cx.B2s(value))
    }

    // Delete a value
    cache.Delete(cx.S2b("key1"))
}
```

# Caveats / Limitations
1. You need to set the eviction count parameter according to usage pattern, it's not a limitation, you can set as 1 or whatever, up to you.
2. Bytes version currently support []byte only as key and value but you can easily convert other types to []byte.
3. Size entry is an estimated size of the cache only. May deviate by 24 - 80 bytes per item or so in actual use. (e.g. assume 1mil entries to have 24mb - 80mb additional overhead)
4. up to 2^63/2 keys for 64 bit system and 2 billion items for 32 bit systems. (not tested on 32bit though, if u need this feature and it doesnt work, drop an issue. will see how to fix for u)

# Roadmap / Todo
- will be changing to Set instead of Put soon
- add more types / generic types, generic version is here, performance is kind of sad but usable. will improve.
https://github.com/cloudxaas/gocache/tree/main/lru
- maybe use a swiss map
- add ttl support

Contributors welcome.
