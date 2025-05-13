package main

import (
	"fmt"
	exam2 "task1/exam-2"
)

func main() {

	// 题目一：只出现一次的数字
	// nums := []int64{1, 2, 3, 2, 1}
	// singleNum := exam1.FindSingleIntSelfMake(nums)
	// fmt.Println(singleNum) // 3
	// nums = []int64{1, 2, 3, 3, 1}
	// singleNum = exam1.FindSingleIntSelfMake(nums)
	// fmt.Println(singleNum) // 2

	// 题目二：回文数
	b1 := exam2.IsPalindromeSelfMake(123)
	b2 := exam2.IsPalindromeSelfMake(121)
	fmt.Println(b1)
	fmt.Println(b2)

	// 题目三：有效的括号
	// str1 := "( [ { 12 } ] )"
	// str2 := "( [ { 12 ) ] )"
	// b1 := exam3.IsValidSelfMake(str1)
	// b2 := exam3.IsValidSelfMake(str2)

	// b3 := exam3.IsValidDoubaoMake(str1)
	// b4 := exam3.IsValidDoubaoMake(str2)

	// fmt.Printf("Self-Exam-Str1: %v\n", b1) // true
	// fmt.Printf("Self-Exam-Str2: %v\n", b2) // false

	// fmt.Printf("Doubao-Exam-Str1: %v\n", b3) // true
	// fmt.Printf("Doubao-Exam-Str2: %v\n", b4) // false
}
