package main

import (
	"fmt"
	exam7 "task3/exam-7"
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
	// post1 := exam6.GetPostWithCommentByUser(2)
	// fmt.Println(post1)
	// post2, err := exam6.GetPostWithMaxComment()
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(*post2)
	// }

	// 题目七
	post := exam7.PostEntity{
		Title:   "标题1",
		Content: "内容1",
		UserId:  2,
	}

	err := exam7.CreatePost(&post)
	if err != nil {
		fmt.Println(err)
	}

	comment1 := exam7.CommentEntity{
		Content: "评论1",
		PostId:  post.Id,
		UserId:  2,
	}
	comment2 := exam7.CommentEntity{
		Content: "评论2",
		PostId:  post.Id,
		UserId:  3,
	}

	err = exam7.AddComment(&comment1)
	if err != nil {
		fmt.Println(err)
	}
	err = exam7.AddComment(&comment2)
	if err != nil {
		fmt.Println(err)
	}

	err = exam7.DeleteComment(&comment1)
	if err != nil {
		fmt.Println(err)
	}
	err = exam7.DeleteComment(&comment2)
	if err != nil {
		fmt.Println(err)
	}
}
