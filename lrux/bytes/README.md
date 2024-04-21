<p align="center">

  <img src="https://github.com/cloudxaas/gocache/assets/104323920/5948a699-64c8-47b8-a5d6-5afedb6a3976" width="40%" height="auto">

    
 <h2 align="center">" Let's <u>Acce</u>lerate <u>LRU</u> " - AccelruX</h2>
</p>

# One of The Fastest Zero Allocation LRU Cache in Golang (for key, value pairs in []byte) - AccelruX (cxlrubytes)

# This X version has higher performance, much lower memory use, at the expense of a bit lower hit ratio
predefine hashing function for hashing keys up to 4 billion keys

Difference with non-x version:

1. non-x uses map[string]int, x uses map[uint32]uint32 (special thanks and credit to phuslu/lru suggestion)

2. non-x does not need a hashing mechanism, you need to predefine the hashing mechanism for x version.
 
3. non-x can use up to available memory and will save all keys as much as possible (may have better hit ratio), x version's hit ratio depends on the hash mechanism used.

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
BenchmarkOtterSet-12                        	 7456851	       149.1 ns/op	      65 B/op	       1 allocs/op
BenchmarkOtterGet-12                        	17926213	        71.64 ns/op	       0 B/op	       0 allocs/op
BenchmarkOtterDelete-12                     	102627896	        11.48 ns/op	       0 B/op	       0 allocs/op
BenchmarkPhusluLRUSet-12                    	66875689	        18.21 ns/op	       0 B/op	       0 allocs/op
BenchmarkPhusluLRUGet-12                    	59487871	        16.95 ns/op	       0 B/op	       0 allocs/op
BenchmarkPhusluLRUDelete-12                 	137954754	         8.624 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesSet-12                   	68247937	        15.28 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesGet-12                   	88466262	        13.78 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesDel-12                   	253530190	         4.740 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesSetParallel-12           	28983519	        47.33 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesGetParallel-12           	25165668	        50.10 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesDelParallel-12           	36730446	        31.12 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesShardedSet-12            	60919240	        19.63 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesShardedGet-12            	100000000	        18.39 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesShardedDel-12            	172000660	         7.133 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesShardedSetParallel-12    	89775692	        12.81 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesShardedGetParallel-12    	65157898	        18.55 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesShardedDelParallel-12    	189985650	         6.418 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/cloudxaas/gocache/lrux/bytes	33.491s
```
reference :

https://github.com/maypok86/otter

https://github.com/phuslu/lru


These benchmarks illustrate the efficiency and speed of the cache, which is designed to operate with zero memory allocations during runtime operations, contributing to its high performance.

## Usage

To use this cache, check the examples folder included, you can configure your own hash function, use xxh3 if you want faster hashing for larger key values > 24 bytes.

# Roadmap / Todo
- add more types / generic types, generic version is here, performance is kind of sad but usable. will improve.
https://github.com/cloudxaas/gocache/tree/main/lru
- maybe use a swiss map
- add ttl support

Contributors welcome.
