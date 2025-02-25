package function

// Set类型定义
type set struct {
	// struct为结构体类型的变量
	m map[any]struct{}
}

// set类型数据结构的初始化操作，在声明的同时可以选择传入或者不传入进去。声明Map切片的时候，Key可以为任意类型的数据，用空接口来实现即可。Value的话按照上面的分析，用空结构体即可
func Set(items ...any) *set {
	// 获取set的地址
	s := &set{}
	// 声明map类型的数据结构
	s.m = make(map[any]struct{})
	s.Add(items...)
	return s
}

// 添加,简化操作可以添加不定个数的元素进入到set中，用变长参数的特性来实现这个需求即可，因为Map不允许Key值相同，所以不必有排重操作。同时将Value数值指定为空结构体类型。
func (s *set) Add(items ...any) {
	for _, item := range items {
		s.m[item] = struct{}{}
	}
}

// 包含,Contains操作其实就是查询操作，看看有没有对应的Item存在，可以利用Map的特性来实现，但是由于不需要Value的数值，所以可以用 _,ok来达到目的：
func (s *set) Contains(item any) bool {
	_, ok := s.m[item]
	return ok
}

// 获取set长度很简单，只需要获取底层实现的Map的长度即可：
func (s *set) Size() int {
	return len(s.m)
}

// 清除操作的话，可以通过重新初始化set来实现，如下即为实现过程：
func (s *set) Clear() {
	s.m = make(map[any]struct{})
}

// 判断两个set是否相等，可以通过循环遍历来实现，即将A中的每一个元素，查询在B中是否存在，只要有一个不存在，A和B就不相等，
func (s *set) Equal(other *set) bool {
	// 如果两者Size不相等，就不用比较了
	if s.Size() != other.Size() {
		return false
	}
	// 迭代查询遍历
	for key := range s.m {
		// 只要有一个不存在就返回false
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

// 判断A是不是B的子集，也是循环遍历的过程，具体分析在上面已经讲述过，实现方式如下所示：
func (s *set) IsSubset(other *set) bool {
	// s的size长于other，不用说了
	if s.Size() > other.Size() {
		return false
	}
	// 迭代遍历
	for key := range s.m {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

// 删除元素
func (s *set) Remove(item any) {
	delete(s.m, item)
}

// 返回集合所有元素
func (s *set) Members() []any {
	var members = []any{}
	// 迭代查询遍历
	for key := range s.m {
		members = append(members, key)
	}
	return members
}
