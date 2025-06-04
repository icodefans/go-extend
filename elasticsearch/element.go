package elasticsearch

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// 解析es标签
func parseESTag(tag string) map[string]string {
	result := make(map[string]string)
	parts := strings.Split(tag, ",")

	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		} else if len(kv) == 1 {
			result["type"] = kv[0]
		}
	}

	return result
}

// 生成Elasticsearch映射
func GenerateMapping(s interface{}) (map[string]interface{}, error) {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", t.Kind())
	}

	properties := make(map[string]interface{})

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 获取json标签作为字段名
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		fieldName := strings.Split(jsonTag, ",")[0]

		// 获取es标签
		esTag := field.Tag.Get("es")
		if esTag == "" {
			continue
		}

		// 解析es标签
		esAttrs := parseESTag(esTag)
		if len(esAttrs) == 0 {
			continue
		}

		// 构建字段映射
		fieldMapping := make(map[string]interface{})
		for k, v := range esAttrs {
			// 处理嵌套属性
			if strings.Contains(v, "{") {
				var nestedValue interface{}
				if err := json.Unmarshal([]byte(v), &nestedValue); err == nil {
					fieldMapping[k] = nestedValue
					continue
				}
			}
			fieldMapping[k] = v
		}

		properties[fieldName] = fieldMapping
	}

	return map[string]interface{}{
		"properties": properties,
	}, nil
}
