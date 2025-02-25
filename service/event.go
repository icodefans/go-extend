package service

import (
	"reflect"
)

// 事件参数
type EventParam struct {
	Key     string  `json:"key"`
	Path    string  `json:"path"`
	Level   string  `json:"level"`
	Error   int     `json:"error"`
	Message *string `json:"message"`
	Label   string  `json:"lable"`
	Data    []any   `json:"data"`
}

// 参数提取
func (event *EventParam) Param(params ...any) {
	for _, arg := range event.Data {
		valueOf := reflect.ValueOf(arg)
		a := valueOf.Kind()
		for _, param := range params {
			valueOf := reflect.ValueOf(param)
			b := valueOf.Kind()
			if a == b {
				param = arg
			}
		}
	}
}

// 事件处理函数
type EventHandlerFunc func(event *EventParam) *Result

type eventHandlerList []EventHandlerFunc

// 事件上下文配置（支持多个事件处理函数）
var eventConfig = map[string]eventHandlerList{
	// "github.com/icodefans/go-extend/service.Error":   {debug,test},
	// "github.com/icodefans/go-extend/service.Success": {debug,test},
}

// 事件触发
func EventTrigger(key, path string, error int, message *string, label, level string, data ...any) {
	funcList, ok := eventConfig[key]
	if !ok {
		return
	}
	for _, f := range funcList {
		f(&EventParam{
			Key:     key,
			Path:    path,
			Error:   error,
			Level:   level,
			Message: message,
			Label:   label,
			Data:    data,
		})
	}
}

// 事件监听
func EventListen(key string, f EventHandlerFunc) {
	eventConfig[key] = append(eventConfig[key], f)
}

// 事件移除
func EventRemove(key string) {
	delete(eventConfig, key)
}
