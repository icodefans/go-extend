// 服务结果返回
package service

import (
	"encoding/json"
	"fmt"
	"runtime"
)

// 结果结构
type Result struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
	Label   string `json:"-"`
	Level   string `json:"-"` // 日志级别
	Key     string `json:"-"` // 事件标识
	Path    string `json:"-"` // 原始方法路径
	Data    any    `json:"data"`
}

// 结构转JONS字符串
func (rs Result) String() []byte {
	jsonStr, _ := json.Marshal(rs)
	return jsonStr
}

// 结构转JONS字符串
func (rs Result) Byte() []byte {
	jsonStr, _ := json.Marshal(rs)
	return jsonStr
}

// 获取错误信息
func (rs Result) GetError() error {
	return fmt.Errorf(rs.Message)
}

// 结果数据结构
type ResultData struct {
	Id uint64 `json:"id,string"`
}

// 结果记录统计
type ResultDataCount struct {
	Count int64 `json:"count"`
}

// 锁定标识
func Key(skip int) string {
	pc, _, _, _ := runtime.Caller(skip)
	return runtime.FuncForPC(pc).Name()
}

// 失败返回
func Error(error int, message string, data ...any) *Result {
	var result = Result{
		Error:   error,
		Message: message,
		Label:   message,
		Key:     Key(1), // 事件标识
		Path:    Key(2), // 原始方法路径
	}
	if len(data) > 0 {
		result.Data = data[0]
		result.Level, _ = data[len(data)-1].(string)
	}
	EventTrigger(result.Key, result.Path, error, &result.Message, result.Label, result.Level, data...)  // trace
	EventTrigger(result.Path, result.Path, error, &result.Message, result.Label, result.Level, data...) // event
	return &result
}

// 成功返回
func Success(message string, data ...any) *Result {
	var result = Result{
		Error:   0,
		Message: message,
		Label:   message,
		Key:     Key(1), // 事件标识
		Path:    Key(2), // 原始方法路径
	}
	if len(data) > 0 {
		result.Data = data[0]
		result.Level, _ = data[len(data)-1].(string)
	}
	EventTrigger(result.Key, result.Path, 0, &result.Message, result.Label, result.Level, data...)  // trace
	EventTrigger(result.Path, result.Path, 0, &result.Message, result.Label, result.Level, data...) // event
	return &result
}
