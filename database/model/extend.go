package model

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
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
	if val, ok := value.([]byte); !ok {
		return nil
	} else {
		decoder := json.NewDecoder(strings.NewReader(string(val)))
		decoder.UseNumber()
		return decoder.Decode(&t)
		// return json.Unmarshal(value.([]byte), &t)
	}
}
