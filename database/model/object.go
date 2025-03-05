package model

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
)

// JSON 对象
type Object map[string]any

// Value 存储数据的时候转换为字符串
func (t Object) Value() (driver.Value, error) {
	if t == nil {
		t = Object{}
	}
	return json.Marshal(t)
}

// Scan 读取数据的时候转换为json(解决uint64类型数据精度丢失问题)
func (t *Object) Scan(value any) error {
	decoder := json.NewDecoder(strings.NewReader(string(value.([]byte))))
	decoder.UseNumber()
	return decoder.Decode(&t)
	// return json.Unmarshal(value.([]byte), &t)
}
