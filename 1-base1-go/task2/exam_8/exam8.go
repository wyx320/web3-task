package exam8

import "fmt"

/*
题目八：
实现一个带有缓冲的通道，
生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
考察点：
通道的缓冲机制。
*/

func Producer(ch chan<- int) {
	for i := 0; i < 100; i++ {
		ch <- i
	}
	close(ch)
}

func Consumer(ch <-chan int) {
	for i := range ch {
		fmt.Println(i)
	}
}
