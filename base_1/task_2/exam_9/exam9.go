package exam9

import "sync"

/*
题目九：
编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。
启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
考察点：
sync.Mutex 的使用、并发数据安全。
*/

func IncreaseCounter(count *int, goroutineCount int) *int {
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mu.Lock()
				(*count)++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return count
}
