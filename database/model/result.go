package model

import (
	"database/sql/driver"
	"encoding/json"
	"strings"

	"github.com/icodefans/go-extend/service"
)

// JSON 对象
type ResultList []service.Result

// Value 存储数据的时候转换为字符串
func (t ResultList) Value() (driver.Value, error) {
	if t == nil {
		t = ResultList{}
	}
	return json.Marshal(t)
}

// Scan 读取数据的时候转换为json(解决uint64类型数据精度丢失问题)
func (t *ResultList) Scan(value any) error {
	decoder := json.NewDecoder(strings.NewReader(string(value.([]byte))))
	decoder.UseNumber()
	return decoder.Decode(&t)
	// return json.Unmarshal(value.([]byte), &t)
}
