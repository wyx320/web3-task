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

// x = 121
// x = 1221
func IsPalindromeTongyiQwQMake(x int) bool {
	// 负数不是回文数，因为负号的存在导致正序和倒序不同
	if x < 0 {
		return false
	}

	// 特殊情况：个位数都是回文数
	if x >= 0 && x < 10 {
		return true
	}

	// 如果以 0 结尾且不等于 0，则一定不是回文数
	if x%10 == 0 && x != 0 {
		return false
	}

	reversed := 0      // 用于存储反转后的数字
	for x > reversed { //121	x(121)>reversed(0)	x(12)>reversed(1)	x(1)>reversed(12)		//1221	x(1221)>reversed(0)	x(122)>reversed(1)	x(12)>reversed(12)
		reversed = reversed*10 + x%10 //121	reversed=1	12										//1221	reversed=12	12
		x /= 10                       //121	x=12	12											//1221	x=12	12
	}

	// 判断原数字是否等于反转后的一半或全部
	return x == reversed || x == reversed/10 //121	x(1)==reversed(11)||x(1)==reversed(11)/10	//1221	x(1)==reversed(11)||x(1)==reversed(11)/10
}

func SelfRepeat(x int) bool {
	if x < 0 {
		return false
	}
	if x >= 0 && x < 10 {
		return true
	}
	if x%10 == 0 {
		return false
	}

	reversed := 0
	for x > reversed {
		reversed = reversed*10 + x%10
		x /= 10
	}

	return x == reversed || x == reversed/10
}
