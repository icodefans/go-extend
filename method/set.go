package method

import "fmt"

// Set 是一个泛型集合类型
type set[T comparable] struct {
	elements map[T]struct{}
}

// Set 创建一个新的集合
func Set[T comparable](eles ...T) *set[T] {
	s := &set[T]{
		elements: make(map[T]struct{}),
	}
	s.Add(eles...)
	return s
}

// Add 向集合中添加多个元素
func (s *set[T]) Add(elements ...T) {
	for _, e := range elements {
		s.elements[e] = struct{}{}
	}
}

// Remove 从集合中移除元素
func (s *set[T]) Remove(element T) {
	delete(s.elements, element)
}

// Contains 判断集合是否包含指定元素
func (s *set[T]) Contains(element T) bool {
	_, exists := s.elements[element]
	return exists
}

// Size 返回集合中元素的数量
func (s *set[T]) Size() int {
	return len(s.elements)
}

// Clear 清空集合
func (s *set[T]) Clear() {
	s.elements = make(map[T]struct{})
}

// Elements 返回集合中所有元素的切片
func (s *set[T]) Elements() []T {
	elements := make([]T, 0, len(s.elements))
	for e := range s.elements {
		elements = append(elements, e)
	}
	return elements
}

// Union 计算两个集合的并集
func (s *set[T]) Union(other *set[T]) *set[T] {
	result := Set[T]()
	for e := range s.elements {
		result.Add(e)
	}
	for e := range other.elements {
		result.Add(e)
	}
	return result
}

// Intersection 计算两个集合的交集
func (s *set[T]) Intersection(other *set[T]) *set[T] {
	result := Set[T]()
	for e := range s.elements {
		if other.Contains(e) {
			result.Add(e)
		}
	}
	return result
}

// Difference 计算两个集合的差集 (s - other)
func (s *set[T]) Difference(other *set[T]) *set[T] {
	result := Set[T]()
	for e := range s.elements {
		if !other.Contains(e) {
			result.Add(e)
		}
	}
	return result
}

// String 实现Stringer接口，方便打印
func (s *set[T]) String() string {
	return fmt.Sprintf("%v", s.Elements())
}

func main1() {
	// 测试整数集合
	intSet := Set[int]()
	intSet.Add(1)
	intSet.Add(2, 3, 4, 5)
	// aaa := intSet.Elements()
	fmt.Println("整数集合:", intSet)
	fmt.Println("集合大小:", intSet.Size())
	fmt.Println("是否包含3:", intSet.Contains(3))

}
func main2() {
	// 测试整数集合
	intSet := Set[int]()
	intSet.Add(1)
	intSet.Add(2, 3, 4, 5)
	fmt.Println("整数集合:", intSet)
	fmt.Println("集合大小:", intSet.Size())
	fmt.Println("是否包含3:", intSet.Contains(3))

	// 测试字符串集合
	strSet := Set[string]()
	strSet.Add("apple")
	strSet.Add("banana", "cherry")
	fmt.Println("字符串集合:", strSet)

	// 测试集合操作
	setA := Set[int]()
	setA.Add(1, 2, 3, 4)

	setB := Set[int]()
	setB.Add(3, 4, 5, 6)

	fmt.Println("setA:", setA)
	fmt.Println("setB:", setB)
	fmt.Println("并集:", setA.Union(setB))
	fmt.Println("交集:", setA.Intersection(setB))
	fmt.Println("差集 (A-B):", setA.Difference(setB))
}

func main3() {
	// 测试整数集合
	intSet := Set(0, 2, 3)
	intSet.Add(1)
	intSet.Add(2, 3, 4, 5)
	// aaa := intSet.Elements()
	fmt.Println("整数集合:", intSet)
	fmt.Println("集合大小:", intSet.Size())
	fmt.Println("是否包含3:", intSet.Contains(3))

}
