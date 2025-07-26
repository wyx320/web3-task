package exam6

import "strconv"

/*
题目六：加一
难度：简单
考察：数组操作、进位处理
题目：给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一
链接：https://leetcode-cn.com/problems/plus-one/

*/

func GetIntBySliceSelfMake(nums []int) (resNums []int) {
	intRes := 0
	for i := 0; i < len(nums); i++ {
		intRes = intRes*10 + nums[i]
	}
	intRes++

	str := strconv.Itoa(intRes)
	for _, value := range str {
		resNums = append(resNums, int(value-'0'))
	}
	return resNums
}

func GetIntBySliceTongyiQwQMake(nums []int) []int {
	for i := len(nums) - 1; i >= 0; i-- {
		if nums[i] < 9 {
			nums[i]++
			return nums
		}
		nums[i] = 0
	}
	return append([]int{1}, nums...)
}
