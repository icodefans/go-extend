package model

import (
	"database/sql/driver"
	"encoding/json"
)

// 扩展数据
type Extend map[string]string

// Value 存储数据的时候转换为字符串
func (t Extend) Value() (driver.Value, error) {
	if t == nil {
		t = Extend{}
	}
	return json.Marshal(t)
}

// Scan 读取数据的时候转换为json
func (t *Extend) Scan(value any) error {
	return json.Unmarshal(value.([]byte), &t)
}
