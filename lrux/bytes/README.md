<p align="center">

  <img src="https://github.com/cloudxaas/gocache/assets/104323920/5948a699-64c8-47b8-a5d6-5afedb6a3976" width="40%" height="auto">

    
 <h2 align="center">" Let's <u>Acce</u>lerate <u>LRU</u> " - AccelruX</h2>
</p>

# One of The Fastest Zero Allocation LRU Cache in Golang (for key, value pairs in []byte) - AccelruX (cxlrubytes)

# This X version has higher performance, much lower memory use, at the expense of maybe a bit lower hit ratio
predefine hashing function for hashing keys up to 4 billion keys (can make this to 2^64 if there are requests for it)

Difference with non-x version:

1. non-x uses map[string]int, x uses map[uint32]uint32 (special thanks and credit to phuslu/lru suggestion).

2. non-x does not need a hashing mechanism (so no collision issues), you need to predefine the hashing mechanism for x version.
 
3. non-x can use up to available memory and will save all keys as much as possible (may have better hit ratio), x version's hit ratio depends on the hash mechanism used, memory can be more efficient but hit ratio may suffer at collision of hash

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

## NOTE : AccelruX capacity is set in terms of MEMORY SIZE LIMIT and not NUMBER OF ITEMS LIMIT.

#### Note : Benchmark results (100kb cache with 1024b key and value to force 1 item eviction here, can set batch eviction higher at your own discretion):
**These benchmarks are for reference only, the memory used is far lesser than most of the rest used AND has higher eviction than the rest (only in this benchmark, for production use, please set it more than 100kb cache please), resulting in much lower hit ratio because most of them are using 10000*(1000+8)bytes ~ 10,080,000 bytes when cxlruxbytes is only using 100,000bytes, 100x lesser memory, for a fairer comparison, use 10mb setting for capacity (but this will NOT result in many evictions, giving it much much higher hit ratio. do test it to verify claims.
**

```
go test -bench=. -benchmem
goos: linux
goarch: amd64
pkg: github.com/cloudxaas/gocache/lrux/bytes
cpu: AMD Ryzen 5 7640HS w/ Radeon 760M Graphics     
BenchmarkHashicorpLRUSet-12                 	59776335	        18.68 ns/op	       0 B/op	       0 allocs/op
BenchmarkHashicorpLRUGet-12                 	69722470	        17.15 ns/op	       0 B/op	       0 allocs/op
BenchmarkHashicorpLRURemove-12              	146819274	         8.295 ns/op	       0 B/op	       0 allocs/op
BenchmarkGoFreeLRUSet-12                    	34314114	        30.47 ns/op	       0 B/op	       0 allocs/op
BenchmarkGoFreeLRUGet-12                    	34552854	        35.80 ns/op	       0 B/op	       0 allocs/op
BenchmarkGoFreeLRURemove-12                 	136824891	         8.748 ns/op	       0 B/op	       0 allocs/op
BenchmarkOtterSet-12*                        	 7792684	       149.6 ns/op	      65 B/op	       1 allocs/op
BenchmarkOtterGet-12                        	17524428	        67.42 ns/op	       0 B/op	       0 allocs/op
BenchmarkOtterDelete-12                     	101478790	        11.86 ns/op	       0 B/op	       0 allocs/op
BenchmarkPhusluLRUSet-12                    	63850137	        19.26 ns/op	       0 B/op	       0 allocs/op
BenchmarkPhusluLRUGet-12                    	72362446	        16.78 ns/op	       0 B/op	       0 allocs/op
BenchmarkPhusluLRUDelete-12                 	130892822	         8.717 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesSet-12                   	77462673	        15.09 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesGet-12                   	88257292	        14.09 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesDel-12                   	253043168	         4.797 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesSetParallel-12           	22535736	        46.59 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesGetParallel-12           	23807607	        49.34 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesDelParallel-12           	37575604	        30.97 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesShardedSet-12            	60373840	        19.40 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesShardedGet-12            	100000000	        12.17 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesShardedDel-12            	168262658	         7.127 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesShardedSetParallel-12    	72022780	        18.45 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesShardedGetParallel-12    	68527062	        18.41 ns/op	       0 B/op	       0 allocs/op
BenchmarkCXLRUBytesShardedDelParallel-12    	185071548	         6.563 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/cloudxaas/gocache/lrux/bytes	58.903s

```
reference :

https://github.com/maypok86/otter

https://github.com/phuslu/lru

https://github.com/hashicorp/golang-lru

https://github.com/elastic/go-freelru

These benchmarks results vary, please adjust parameters to use for your own use case. 
*otter has 1 byte alloc / op, in these zero alloc benchmark

## NOTE : AccelruX capacity is set in terms of MEMORY SIZE LIMIT and not NUMBER OF ITEMS LIMIT.

## Usage

To use this cache, check the examples folder included, you can configure your own hash function, use xxh3 if you want faster hashing for larger key values > 24 bytes.

# Roadmap / Todo
- add more types / generic types, generic version is here, performance is kind of sad but usable. will improve.
https://github.com/cloudxaas/gocache/tree/main/lru
- add ttl support

Contributors welcome.
