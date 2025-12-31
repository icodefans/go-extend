package function

// ReverseSlice 原地反转切片（任意类型）
func ReverseSlice[T any](s []T) {
	// 双指针：i 从头部开始，j 从尾部开始
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		// 交换首尾元素
		s[i], s[j] = s[j], s[i]
	}
}

// ReverseSliceNew 返回反转后的新切片（不修改原切片）
func ReverseSliceNew[T any](s []T) []T {
	// 创建新切片，容量和长度与原切片一致
	newSlice := make([]T, len(s))
	// 从后往前遍历原切片，赋值给新切片
	for i := range s {
		newSlice[i] = s[len(s)-1-i]
	}
	return newSlice
}
