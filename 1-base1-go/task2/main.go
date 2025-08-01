package main

import (
	"fmt"
	exam10 "task2/exam-10"
)

func main() {

	// 题目一
	// i := 7
	// exam1.PtrMethod(&i)
	// fmt.Println(i)

	// 题目二
	// nums := []int{1, 3, 5, 7, 9}
	// exam2.IntSlicePreElementMultiplyTwo(&nums)
	// fmt.Println(nums)

	// 题目三
	// var wg sync.WaitGroup
	// wg.Add(2)
	// exam3.DoubleCoroutine(&wg)
	// wg.Wait()

	// 题目四
	// exam4.TaskSceduler([](func() string){
	// 	exam4.CreateTask("A", time.Second*1),
	// 	exam4.CreateTask("B", time.Second*2),
	// 	exam4.CreateTask("C", time.Second*3),
	// 	exam4.CreateTask("D", time.Second*4),
	// })
	// exam4.TaskSceduler([]exam4.Task{
	// 	exam4.CreateTask("A", time.Second*1),
	// 	exam4.CreateTask("B", time.Second*2),
	// 	exam4.CreateTask("C", time.Second*3),
	// 	exam4.CreateTask("D", time.Second*4),
	// })

	// 题目五
	// rect := exam5.Rectange{Length: 10, Width: 2}
	// fmt.Println(rect.Area())      // 20
	// fmt.Println(rect.Perimeter()) // 24
	// cir := exam5.Circle{Radius: 5}
	// fmt.Println(cir.Area())      // 78.54
	// fmt.Println(cir.Perimeter()) // 31.42

	// 题目六
	// emp := exam6.Employee{EmployeeID: 1, Person: exam6.Person{Name: "张三", Age: 18}}
	// emp.PrintInfo()

	// 题目七
	// ch := make(chan int, 10)
	// go exam7.Producer(ch)
	// exam7.Consumer(ch)

	// 题目八
	// ch := make(chan int, 100)
	// go exam8.Producer(ch)
	// exam8.Consumer(ch)

	// 题目九
	// count := 0
	// exam9.IncreaseCounter(&count, 10)
	// fmt.Println(count)

	// 题目十
	var count uint64 = 0
	exam10.IncreaseCounter(&count, 10)
	fmt.Println(count)
}
