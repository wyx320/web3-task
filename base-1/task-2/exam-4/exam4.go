package exam4

import (
	"fmt"
	"sync"
	"time"
)

/*
题目四：
设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
考察点：
协程原理、并发任务调度。
*/

// type Task func() string

// func TaskSceduler(tasks []Task) {
// 	var wg sync.WaitGroup
// 	for _, task := range tasks {
// 		wg.Add(1)
// 		go func(t Task) {
// 			defer wg.Done()

// 			startTime := time.Now()
// 			taskName := task()
// 			endTime := time.Now()

// 			fmt.Printf("TaskName:%s DurationTime:%v\n", taskName, endTime.Sub(startTime))
// 		}(task)
// 	}
// 	wg.Wait()
// }

// func CreateTask(taskName string, durationTime time.Duration) Task {
// 	return func() string {
// 		time.Sleep(durationTime)
// 		return taskName
// 	}
// }

// 熟能生巧

type Task func() string

func TaskSceduler(tasks []Task) {
	results := make([]string, len(tasks))

	var wg sync.WaitGroup
	for index, task := range tasks {
		wg.Add(1)

		go func(index int, task Task) {
			defer wg.Done()

			startTime := time.Now()
			result := task()
			endTime := time.Now()

			executionTime := endTime.Sub(startTime)
			results[index] = fmt.Sprintf("%v (Execution Time: %v)", result, executionTime)
		}(index, task)
	}
	wg.Wait()

	for _, result := range results {
		fmt.Println(result)
	}
}

func CreateTask(taskName string, durationTime time.Duration) Task {
	return func() string {
		time.Sleep(durationTime)
		return fmt.Sprintf("%v Execution Completed. ", taskName)
	}
}
