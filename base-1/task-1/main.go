package main

import (
	"fmt"
	exam8 "task1/exam-8"
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
	// b1 := exam2.IsPalindromeSelfMake(123)
	// b2 := exam2.IsPalindromeSelfMake(121)
	// fmt.Println(b1)
	// fmt.Println(b2)
	// b3 := exam2.IsPalindromeTongyiQwQMake(123)
	// b4 := exam2.IsPalindromeTongyiQwQMake(121)
	// fmt.Println(b3)
	// fmt.Println(b4)

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

	// 题目四：查找字符串数组中的最长公共前缀
	// strs := []string{"flower", "flow", "flight"}
	// res := exam4.FindPublicPrefixSelfMake(strs)
	// fmt.Println(res)

	// 题目五：删除排序数组中的重复项
	// nums := []int{1, 2, 2, 3, 4, 5, 5, 6}   // 输入数组
	// expectedNums := []int{1, 2, 3, 4, 5, 6} // 长度正确的期望答案
	// k := exam5.RemoveDuplicates(&nums)      // 调用
	// if k != len(expectedNums) {
	// 	fmt.Println("唯一元素个数不符合预期")
	// 	return
	// }
	// for i := 0; i < k; i++ {
	// 	if nums[i] != expectedNums[i] {
	// 		fmt.Println("处理后的数组不符合预期")
	// 		return
	// 	}
	// }
	// fmt.Printf("%#v\n", nums)
	// fmt.Printf("唯一元素长度：%v", k)

	// 题目六：加一
	// nums := []int{4, 3, 2, 1}
	// fmt.Println(exam6.GetIntBySliceSelfMake(nums))
	// nums = []int{9}
	// fmt.Println(exam6.GetIntBySliceSelfMake(nums))
	// nums = []int{4, 3, 2, 1}
	// fmt.Println(exam6.GetIntBySliceTongyiQwQMake(nums))
	// nums = []int{9}
	// fmt.Println(exam6.GetIntBySliceTongyiQwQMake(nums))

	// 题目七：26. 删除有序数组中的重复项：
	// elementCount := exam7.RemoveDuplicates([]int{1, 2, 2, 3, 4, 5, 5, 6})
	// fmt.Println(elementCount)

	// 题目八：56. 合并区间
	var res = exam8.ConcatRangeSelfMake([][]float64{{2, 5}, {1, 3}})
	fmt.Printf("%#v\n", res)
	res = exam8.ConcatRangeTongyiQwQMake([][]float64{{2, 5}, {1, 3}})
	fmt.Printf("%#v", res)
}
