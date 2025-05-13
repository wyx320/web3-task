package exam2

/*
题目二：有效的括号：

给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。
有效字符串需满足：
左括号必须用相同类型的右括号闭合。
左括号必须以正确的顺序闭合。
每个右括号都有一个对应的相同类型的左括号。

链接：https://leetcode-cn.com/problems/valid-parentheses/
*/

// self-make
func IsValidSelfMake(str string) (validSucceed bool) {

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

// doubao-make
func IsValidDoubaoMake(str string) bool {
	stack := []rune{}
	for _, char := range str {
		switch char {
		case '(', '{', '[':
			stack = append(stack, char)
		case ')':
			if len(stack) == 0 || stack[len(stack)-1] != '(' {
				return false
			}
			stack = stack[:len(stack)-1]
		case '}':
			if len(stack) == 0 || stack[len(stack)-1] != '{' {
				return false
			}
			stack = stack[:len(stack)-1]
		case ']':
			if len(stack) == 0 || stack[len(stack)-1] != '[' {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}
