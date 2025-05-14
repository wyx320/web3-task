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

func FindPublicPrefixDelfMake(strSlice []string) string {
	resSlice := []rune{}

	tempRune := []rune{}

	index := 0
	for _, str := range strSlice {

		if len(str) == 0 {
			break
		}

		if index > len(str)-1 {
			break
		}

		tempRune = append(tempRune, []rune(str)[index])

		preLetter := tempRune[0]
		isAllSampleLetter := true
		for _, value := range tempRune[1:] {
			if preLetter != value {
				isAllSampleLetter = false
				break
			}
		}
		if isAllSampleLetter {
			resSlice = append(resSlice, preLetter)
		}

		index++
	}
	return string(resSlice)
}
