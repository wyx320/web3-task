package exam2

/*
题目二：
实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
考察点：
指针运算、切片操作。
*/

func IntSlicePreElementMultiplyTwo(slicePtr *[]int) {
	slice := *slicePtr
	for index := range slice {
		slice[index] *= 2
	}
}
