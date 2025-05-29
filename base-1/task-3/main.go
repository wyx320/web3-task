package main

import (
	"fmt"
	exam6 "task3/exam-6"
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
	// emps := exam3.QueryByDepartment("技术部")
	// fmt.Println(emps)
	// emp := exam3.QueryMaxSalary()
	// fmt.Println(emp)

	// 题目四
	// books := exam4.QueryByPrice(50)
	// fmt.Println(books)

	// 题目五
	// exam5.Test()

	// 题目六
	post1 := exam6.GetPostWithCommentByUser(2)
	fmt.Println(post1)
	post2, err := exam6.GetPostWithMaxComment()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(*post2)
	}
}
