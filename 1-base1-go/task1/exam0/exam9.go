package exam9

/*
题目九：两数之和
考察：数组遍历、map使用
题目：给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数
链接：https://leetcode-cn.com/problems/two-sum/

给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 target  的那 两个 整数，并返回它们的数组下标。
你可以假设每种输入只会对应一个答案，并且你不能使用两次相同的元素。
你可以按任意顺序返回答案。

示例 1：
输入：nums = [2,7,11,15], target = 9
输出：[0,1]
解释：因为 nums[0] + nums[1] == 9 ，返回 [0, 1] 。

示例 2：
输入：nums = [3,2,4], target = 6
输出：[1,2]

示例 3：
输入：nums = [3,3], target = 6
输出：[0,1]
*/

func SumTargetSelfMake(digits []int, target int) []int {
	for i := 0; i < len(digits)-1; i++ {
		for j := i + 1; j < len(digits); j++ {
			if digits[i]+digits[j] == target {
				return []int{i, j}
			}
		}
	}

	return []int{}
}

func SumTargetTongyiQwQMake(digits []int, target int) []int {
	m := map[int]int{}
	for index, value := range digits {
		// 计算补数
		complement := target - value

		// 检查补数是否在 map 中
		if i, ok := m[complement]; ok {
			// 找到答案，返回两个数的索引
			return []int{i, index}
		}

		// 将当前数字及索引存入 map
		m[value] = index
	}

	// 如果没有找到答案，返回空切片
	return nil
}
