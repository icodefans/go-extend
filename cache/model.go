package cache

import (
	"fmt"
	"runtime"

	"github.com/icodefans/go-extend/database/redis"
)

type Mode uint8

// 操作方式
const (
	INIT Mode = iota
	READ
	DELETE
	SET
)

// 缓存结果
type Result struct {
	Data   any
	Key    string
	Expire uint32
	// Close  bool
	Redis redis.Redis
	Error error
}

// 分页参数
type Page struct {
	Number int
	Limit  int
}

// 缓存函数类型
type DataHandlerFunc func(mode Mode, args ...any) *Result

// 分页缓存函数类型
type PageHandlerFunc func(mode Mode, subKey string, page *Page, args ...any) *Result

type PagesHandlerFunc func(mode Mode, page *Page, args ...any) *Result

// 缓存标识
func Key(mark *string) string {
	pc, _, _, _ := runtime.Caller(1)
	key := runtime.FuncForPC(pc).Name()
	if mark != nil {
		key = fmt.Sprint(key, "_", *mark)
	}
	return key
}
