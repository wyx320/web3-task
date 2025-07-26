package exam7

/*
题目七：26. 删除有序数组中的重复项：
给你一个有序数组 nums ，请你原地删除重复出现的元素，使每个元素只出现一次，返回删除后数组的新长度。
不要使用额外的数组空间，你必须在原地修改输入数组并在使用 O(1) 额外空间的条件下完成。
可以使用双指针法，一个慢指针 i 用于记录不重复元素的位置，一个快指针 j 用于遍历数组，当 nums[i] 与 nums[j] 不相等时，将 nums[j] 赋值给 nums[i + 1]，并将 i 后移一位。
*/
func RemoveDuplicates(digits []int) int {
	i := 0
	for j := 1; j < len(digits); j++ {
		if digits[i] != digits[j] {
			i++
			digits[i] = digits[j]
		}
	}
	return i + 1
}
