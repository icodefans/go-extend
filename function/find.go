package function

import (
	"fmt"
)

// 列表结构中查找字段
// 使用自定义查找函数
func FindInSlice[T any](slice []T, predicate func(T) bool) (*T, bool) {
	for i, item := range slice {
		if predicate(item) {
			return &slice[i], true
		}
	}
	return nil, false
}

// 使用示例
func main() {
	type Person struct {
		ID   int
		Name string
		Age  int
	}
	people := []Person{}
	person, found := FindInSlice(people, func(p Person) bool {
		return p.Name == "John"
	})
	fmt.Println(person, found)
}
