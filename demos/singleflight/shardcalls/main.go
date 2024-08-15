package main

import (
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type ShardCalls struct {
	shards    []singleflight.Group
	shardMask uint64
}

func NewShardCalls(shardCount int) *ShardCalls {
	shardCount = nextPowerOfTwo(shardCount)
	return &ShardCalls{
		shards:    make([]singleflight.Group, shardCount),
		shardMask: uint64(shardCount - 1),
	}
}

func (sc *ShardCalls) Do(key string, fn func() (interface{}, error)) (interface{}, error, bool) {
	shard := sc.getShard(key)
	return sc.shards[shard].Do(key, fn)
}

func (sc *ShardCalls) getShard(key string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(key))
	return h.Sum64() & sc.shardMask
}

func nextPowerOfTwo(v int) int {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}

func main() {
	sc := NewShardCalls(4) // 創建 4 個分片

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i%10) // 使用 10 個不同的 key
			result, err, shared := sc.Do(key, func() (interface{}, error) {
				time.Sleep(100 * time.Millisecond) // 模擬耗時操作
				return fmt.Sprintf("Result for %s", key), nil
			})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Result: %v, Shared: %v\n", result, shared)
			}
		}(i)
	}

	wg.Wait()
}
