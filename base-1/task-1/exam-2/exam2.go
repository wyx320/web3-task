package exam2

import "strconv"

/*
题目一：回文数

给你一个整数 x ，如果 x 是一个回文整数，返回 true ；否则，返回 false 。
回文数是指正序（从左向右）和倒序（从右向左）读都是一样的整数。
例如，121 是回文，而 123 不是。
*/

func IsPalindromeSelfMake(num int64) bool {
	numStr := strconv.FormatInt(num, 10)
	numStrLen := len(numStr)
	stack := numStr[:numStrLen]
	for i := numStrLen - 1; i > numStrLen/2; i-- {
		if stack[numStrLen-1-i] != numStr[i] {
			return false
		}
	}
	return true
}

func isPalindrome(x int) bool {
	if x < 0 || (x%10 == 0 && x != 0) {
		return false
	}
	reversed := 0
	for x > reversed {
		digit := x % 10
		reversed = reversed*10 + digit
		x /= 10
	}
	return x == reversed || x == reversed/10
}
