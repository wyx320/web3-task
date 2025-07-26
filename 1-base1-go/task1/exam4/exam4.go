package exam4

/*
题目四：查找字符串数组中的最长公共前缀

编写一个函数来查找字符串数组中的最长公共前缀。
如果不存在公共前缀，返回空字符串 ""。

示例 1：
输入：strs = ["flower","flow","flight"]
输出："fl"

示例 2：
输入：strs = ["dog","racecar","car"]
输出：""
解释：输入不存在公共前缀。
*/

func FindPublicPrefixSelfMake(strSlice []string) string {
	if len(strSlice) == 0 {
		return "" // 如果数组为空，则返回空字符串
	}

	// 假设第一个字符串是公共前缀
	prefix := strSlice[0]

	for i := 1; i < len(strSlice); i++ {
		for !startWith(strSlice[i], prefix) {
			// 缩短前缀
			prefix = prefix[:len(prefix)-1]

			if prefix == "" {
				return "" // 如果前缀为空，直接返回空字符串
			}
		}
	}

	return prefix
}

// startWith 判断字符串 str 是否以 prefix 开头
func startWith(str string, prefix string) bool {
	if len(str) < len(prefix) {
		return false
	}
	return str[:len(prefix)] == prefix
}
