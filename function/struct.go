package function

import (
	"github.com/mitchellh/mapstructure"
)

// 结构转字典
func StructToMap(data any) (content map[string]any) {
	_ = mapstructure.Decode(data, &content)
	return content
}
