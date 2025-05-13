package main

import (
	"fmt"
	exam2 "task1/exam-2"
)

func main() {

	// 题目一：回文数

	// 题目二：有效的括号
	str1 := "( [ { 12 } ] )"
	str2 := "( [ { 12 ) ] )"
	b1 := exam2.IsValidSelfMake(str1)
	b2 := exam2.IsValidSelfMake(str2)

	b3 := exam2.IsValidDoubaoMake(str1)
	b4 := exam2.IsValidDoubaoMake(str2)

	fmt.Printf("Self-Exam-Str1: %v\n", b1) // true
	fmt.Printf("Self-Exam-Str2: %v\n", b2) // false

	fmt.Printf("Doubao-Exam-Str1: %v\n", b3) // true
	fmt.Printf("Doubao-Exam-Str2: %v\n", b4) // false
}
