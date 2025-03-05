package function

import (
	"reflect"
)

// 移除键值对中为nil的项
func MapFilterNull(data *map[string]any) *map[string]any {
	for key, value := range *data {
		vi := reflect.ValueOf(value)
		if vi.Kind() == reflect.Ptr && vi.IsNil() {
			delete(*data, key)
		}
	}
	return data
}

// 获取字典key值列表
func MapStringKeys(data any) []string {
	var mapKeys []string
	keys := reflect.ValueOf(data).MapKeys()
	for _, key := range keys {
		mapKeys = append(mapKeys, key.String())
	}
	return mapKeys
}
