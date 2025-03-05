package util

// 字符串数组
type ArrayString map[string]string

// 根据值查找键值
func (arr *ArrayString) Search(value string) *string {
	for key, val := range *arr {
		if val == value {
			return &key
		}
	}
	return nil
}
