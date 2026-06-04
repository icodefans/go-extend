package mysql

import (
	"encoding/json"
	"strconv"
)

// FormatWhereVal 处理 [field,op,val] 第三个元素
func FormatWhereVal(arr *[3]any) {
	val := arr[2]
	res := anyToReal(val)
	arr[2] = res
}

func anyToReal(v any) any {
	switch d := v.(type) {
	case json.Number:
		// json.Number优先int
		if i, err := d.Int64(); err == nil {
			return int(i)
		}
		f, _ := d.Float64()
		return f
	case string:
		// 字符串是数字则转数值
		if n, err := strconv.Atoi(d); err == nil {
			return n
		}
		if f, err := strconv.ParseFloat(d, 64); err == nil {
			return f
		}
		return d
	default:
		return d
	}
}

// ParseSearchWhere 批量格式化全部Where条件
func ParseSearchWhere(where [][3]any) {
	for i := range where {
		FormatWhereVal(&where[i])
	}
}
