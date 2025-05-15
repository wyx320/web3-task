package exam8

import "sort"

/*
题目八：56. 合并区间
要求：
以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
请你合并所有重叠的区间，并返回一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间。
提示：
可以先对区间数组按照区间的起始位置进行排序，然后使用一个切片来存储合并后的区间，遍历排序后的区间数组，将当前区间与切片中最后一个区间进行比较，
如果有重叠，则合并区间；如果没有重叠，则将当前区间添加到切片中。
*/

func ConcatRangeSelfMake(intervals [][]float64) [][]float64 {

	// 冒泡排序
	for i := 0; i < len(intervals); i++ {
		for j := 1; j < len(intervals)-i; j++ {
			if intervals[i][0] > intervals[j][0] || // 左区间从小到大排序
				(intervals[i][0] == intervals[j][0] && intervals[i][1] > intervals[j][1]) { // 左区间相等，右区间从小到大排序
				tempInterval := intervals[i]
				intervals[i] = intervals[j]
				intervals[j] = tempInterval
			}
		}
	}

	intervalResult := [][]float64{}
	// intervalResult = append(intervalResult, intervals[0])
	for i := 0; i < len(intervals); i++ {
		intervalTemp := []float64{}
		if intervals[i][1] > intervals[i+1][0] {
			intervalTemp = []float64{intervals[i][0], intervals[i+1][1]} // 下个元素起 小于 当前元素止 ==》生成临时元素
			i++
		} else {
			intervalTemp = intervals[i] // 下一个元素与当前元素没有重合区间 ==》临时元素=当前元素
		}

		if len(intervalResult) > 0 && intervalResult[len(intervalResult)-1][1] > intervalTemp[0] {
			intervalResult[len(intervalResult)-1][1] = intervalTemp[1] // 切片最后一个元素止 大于 待处理元素起 ==》合并元素
		} else {
			intervalResult = append(intervalResult, intervalTemp) // 临时元素与切片最末尾元素没有重合区间 直接把元素放追加到切片
		}
	}

	return intervalResult
}

func ConcatRangeTongyiQwQMake(intervals [][]float64) [][]float64 {
	if len(intervals) < 2 {
		return intervals
	}

	// 使用 sort.Slice 对区间按照起始位置进行排序
	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i][0] == intervals[j][0] {
			return intervals[i][1] < intervals[j][1]
		} else {
			return intervals[i][0] < intervals[j][0]
		}
	})

	intervalResult := [][]float64{}
	intervalResult = append(intervalResult, intervals[0])

	for i := 1; i < len(intervals); i++ {
		current := intervals[i]
		lastMerged := intervalResult[len(intervalResult)-1]

		if lastMerged[1] >= current[0] {
			lastMerged[1] = max(current[1], lastMerged[1])
		} else {
			intervalResult = append(intervalResult, current)
		}
	}

	return intervalResult
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
