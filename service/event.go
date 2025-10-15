package service

import (
	"time"
)

// 事件参数
type EventParam struct {
	Key       string    `json:"key"`
	Path      string    `json:"path"`
	Level     string    `json:"level"`
	Error     int       `json:"error"`
	StartTime time.Time `json:"start_time"`
	Message   *string   `json:"message"`
	Label     string    `json:"lable"`
	Data      []any     `json:"data"`
}

// 事件处理函数
type EventHandlerFunc func(event *EventParam) *Result

type eventHandlerList []EventHandlerFunc

// 事件上下文配置（支持多个事件处理函数）
var eventConfig = map[string]eventHandlerList{
	"github.com/icodefans/go-extend/command.ApiServerStart": {Trace},
}

// 事件触发
func EventTrigger(key, path string, error int, message *string, label, level string, data ...any) {
	funcList, ok := eventConfig[key]
	if !ok {
		return
	}
	for _, f := range funcList {
		f(&EventParam{
			Key:       key,
			Path:      path,
			Error:     error,
			Level:     level,
			StartTime: time.Now(),
			Message:   message,
			Label:     label,
			Data:      data,
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
