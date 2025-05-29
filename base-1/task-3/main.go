package main

import (
	"fmt"
	exam3 "task3/exam-3"
)

func main() {

	// 题目一：

	// // exam1.CreateNewOne()

	// students := exam1.QueryByAge()
	// fmt.Printf("%#v\n", students)

	// exam1.UpdateGradeByName("张三", "四年级")

	// exam1.DeleteLessThanByAge()

	// 题目二
	// exam2.CreateAccount(50)
	// exam2.CreateAccount(2000)
	// exam2.CreateAccount(0)

	// err := exam2.Transfer(1, 3, 100)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("转账成功")
	// }

	// err = exam2.Transfer(2, 3, 100)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("转账成功")
	// }

	// 题目三
	emps := exam3.QueryByDepartment("技术部")
	fmt.Println(emps)
	emp := exam3.QueryMaxSalary()
	fmt.Println(emp)
}
