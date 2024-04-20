<p align="center">

  <img src="https://github.com/cloudxaas/gocache/assets/104323920/5948a699-64c8-47b8-a5d6-5afedb6a3976" width="40%" height="auto">

    
 <h2 align="center">" Let's <u>Acce</u>lerate <u>LRU</u> " - AccelruX</h2>
</p>

# One of The Fastest Zero Allocation LRU Cache in Golang (for key, value pairs in []byte) - AccelruX (cxlrubytes)

# This X version has higher performance, much lower memory use, at the expense of a bit lower hit ratio
predefine hashing function for hashing keys up to 4 billion keys

Welcome to the repository for one of the fastest LRU cache implementations available for Go. This LRU cache is uniquely designed to limit the memory usage directly, rather than by the number of entries. This makes it ideal for applications where the memory footprint is critical, such as in embedded systems or high-performance computing environments where resources are tightly managed.

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

Benchmark results (100kb cache with 1024b key and value to force 1 item eviction here, can set batch eviction higher at your own discretion):
```
go test -bench=. -benchmem
goos: linux
goarch: amd64
pkg: github.com/cloudxaas/gocache/lrux/bytes
cpu: AMD Ryzen 5 7640HS w/ Radeon 760M Graphics     
BenchmarkPhusluLRUSet-12                        43065601                27.89 ns/op            0 B/op          0 allocs/op
BenchmarkPhusluLRUGet-12                        71053752                17.09 ns/op            0 B/op          0 allocs/op
BenchmarkPhusluLRUDelete-12                     131928775                9.061 ns/op           0 B/op          0 allocs/op
BenchmarkCXLRUBytesSet-12                       76250018                15.86 ns/op            0 B/op          0 allocs/op
BenchmarkCXLRUBytesGet-12                       83375414                14.05 ns/op            0 B/op          0 allocs/op
BenchmarkCXLRUBytesDel-12                       251063827                4.825 ns/op           0 B/op          0 allocs/op
BenchmarkCXLRUBytesSetParallel-12               28785394                48.20 ns/op            0 B/op          0 allocs/op
BenchmarkCXLRUBytesGetParallel-12               26214566                46.60 ns/op            0 B/op          0 allocs/op
BenchmarkCXLRUBytesDelParallel-12               34564299                33.40 ns/op            0 B/op          0 allocs/op
BenchmarkCXLRUBytesShardedSet-12                52624578                21.94 ns/op            0 B/op          0 allocs/op
BenchmarkCXLRUBytesShardedGet-12                89751607                12.56 ns/op            0 B/op          0 allocs/op
BenchmarkCXLRUBytesShardedDel-12                166306846                7.186 ns/op           0 B/op          0 allocs/op
BenchmarkCXLRUBytesShardedSetParallel-12        65316200                18.20 ns/op            0 B/op          0 allocs/op
BenchmarkCXLRUBytesShardedGetParallel-12        59256444                21.89 ns/op            0 B/op          0 allocs/op
BenchmarkCXLRUBytesShardedDelParallel-12        175149573                7.826 ns/op           0 B/op          0 allocs/op
PASS
ok      github.com/cloudxaas/gocache/lrux/bytes 28.996s
```
reference : https://github.com/phuslu/lru


These benchmarks illustrate the efficiency and speed of the cache, which is designed to operate with zero memory allocations during runtime operations, contributing to its high performance.

## Usage

To use this cache, check the examples folder included, you can configure your own hash function, use xxh3 if you want faster hashing for larger key values > 24 bytes.

# Roadmap / Todo
- add more types / generic types, generic version is here, performance is kind of sad but usable. will improve.
https://github.com/cloudxaas/gocache/tree/main/lru
- maybe use a swiss map
- add ttl support

Contributors welcome.
