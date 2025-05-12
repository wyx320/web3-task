package make

import "strings"

func main() {

}

/*
给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。

有效字符串需满足：

左括号必须用相同类型的右括号闭合。
左括号必须以正确的顺序闭合。
每个右括号都有一个对应的相同类型的左括号。
*/
func stringValidator(str string) (validSucceed bool) {
	subStrStart := "({["
	subStrEnd := ")}]"
	for _, value := range str {
		validSucceed = true
		if strings.Contains(subStrStart, string(value)) {
			runeStr := []rune(str)
			for i := len(runeStr) - 1; i >= 0; i-- {
				if !strings.Contains(subStrEnd, string(runeStr[i])) {
					validSucceed = false
					return validSucceed
				}
			}
		}
	}
	return validSucceed
}
