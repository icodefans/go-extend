package logger

import (
	"fmt"
	"log"
)

// 标准错误输出
func StdErrAdd(content, label string) {
	// 日志前缀
	if label == "" {
		label = "EXTEND"
	}
	prefix := fmt.Sprintf("[%s] ", label)
	// 创建一个新的日志记录器，将日志输出到标准错误输出（stderr）
	logger := log.New(log.Writer(), prefix, log.Ldate|log.Ltime)

	// 记录一条错误日志
	logger.Println(content)
}
