package main

import "fmt"

func main() {

	// 题目一：回文数

	// 题目二：有效的括号
	b1 := validBrackets("( [ { 12 } ] )")
	b2 := validBrackets("( [ { 12 ) ] )")
	fmt.Println(b1) // true
	fmt.Println(b2) // false
}

/*
题目一：回文数

给你一个整数 x ，如果 x 是一个回文整数，返回 true ；否则，返回 false 。
回文数是指正序（从左向右）和倒序（从右向左）读都是一样的整数。
例如，121 是回文，而 123 不是。
*/

/*
题目二：有效的括号：

给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。
有效字符串需满足：
左括号必须用相同类型的右括号闭合。
左括号必须以正确的顺序闭合。
每个右括号都有一个对应的相同类型的左括号。

链接：https://leetcode-cn.com/problems/valid-parentheses/
*/
func validBrackets(str string) (validSucceed bool) {

	validSucceed = true // 默认有效

	m := map[rune]rune{
		'(': ')',
		'{': '}',
		'[': ']',
	}
	runeStr := []rune(str)

	leftIndex := 0                 // 这次匹配到左括号时的位置
	rightIndex := len(runeStr) - 1 // 上次匹配到右括号时的位置

	for index, value := range runeStr {

		rightValue, ok := m[value] // 判断是否是左括号

		if ok {
			leftIndex = index // 记录左括号的位置
			for               // 从右往左遍历
			i := rightIndex;  // 右边已遍历的位置，不重复遍历
			i > leftIndex;    // 左边已遍历的位置，不遍历
			i-- {
				// 如果右括号在左边括号右边，而且找到匹配的括号，则本次查找有效
				if runeStr[i] == rightValue {
					rightIndex = i - 1
					break
				}
				if i == leftIndex+1 { // 如果右边没有找到匹配的括号，则该字符串无效
					{
						validSucceed = false
						return
					}
				}
			}
		}
	}

	return validSucceed
}
