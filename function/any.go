package function

import (
	"reflect"
)

// Go语言接口变量的nil判断
func IsNil(v any) bool {
	valueOf := reflect.ValueOf(v)

	k := valueOf.Kind()
	// 反射判断nil时需要注意类型只能是 pointer, channel, func, interface, map, or slice type，使用其他类型会直接 panic。
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return valueOf.IsNil()
	default:
		return v == nil
	}
}
