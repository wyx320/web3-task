package exam5

/*
题目五：删除排序数组中的重复项
难度：简单
考察：数组/切片操作
题目：给定一个排序数组，你需要在原地删除重复出现的元素
链接：https://leetcode-cn.com/problems/remove-duplicates-from-sorted-array/

举例：
source: [1,2,2,3,4,5,5,6]
target: [1,2,3,4,5,6] 6(删除后数组长度) 4(唯一元素个数)

*/

func RemoveDuplicates(numsPtr *[]int) int {
	nums := *numsPtr
	lenNums := len(nums)
	if lenNums == 0 {
		return 0
	}

	k := 0
	for i := 1; i < lenNums; i++ {
		if nums[i] != nums[k] {
			k++
			nums[k] = nums[i]
		}
	}

	*numsPtr = (*numsPtr)[:k+1]
	return k + 1
}
