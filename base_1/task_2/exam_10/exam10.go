package exam10

import (
	"sync"
	"sync/atomic"
)

/*
题目：
使用原子操作（ sync/atomic 包）实现一个无锁的计数器。
启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
考察点：
原子操作、并发数据安全。
*/

func IncreaseCounter(count *uint64, goroutineCount int) {
	var wg sync.WaitGroup
	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddUint64(count, 1)
			}
		}()
	}
	wg.Wait()
}
