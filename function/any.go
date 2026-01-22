package function

import (
	"fmt"
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

// Any变量类型信息打印
func AnyTypePrint(val any) {
	// 获取反射类型对象
	t := reflect.TypeOf(val)
	fmt.Println("变量类型为:", t)
	// 获取反射值对象
	v := reflect.ValueOf(val)
	fmt.Println("变量值为:", v)
	fmt.Println("变量值的类型为:", v.Type())
}
