package model

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

// JSON 对象列表
type Remark []Object

// Value 存储数据的时候转换为字符串
func (t Remark) Value() (driver.Value, error) {
	if t == nil {
		return json.Marshal([]struct{}{})
	}
	t[len(t)-1]["timestamp"] = time.Now().Unix()
	return json.Marshal(t)
}

// Scan 读取数据的时候转换为json(解决uint64类型数据精度丢失问题)
func (t *Remark) Scan(value any) error {
	if _, ok := value.([]byte); !ok {
		return nil
	}
	decoder := json.NewDecoder(strings.NewReader(string(value.([]byte))))
	decoder.UseNumber()
	return decoder.Decode(&t)
	// return json.Unmarshal(value.([]byte), &t)
}
