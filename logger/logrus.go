package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 日志等级字符串关联
var LogrusLevel = map[string]logrus.Level{
	"panic": logrus.PanicLevel,
	"fatal": logrus.FatalLevel,
	"error": logrus.ErrorLevel,
	"warn":  logrus.WarnLevel,
	"info":  logrus.InfoLevel,
	"debug": logrus.DebugLevel,
	"trace": logrus.TraceLevel,
}

// 文件日志对象初始化
func NewRotateLogrus(logName string) *logrus.Logger {
	if logName == "" {
		fmt.Println("logName not null")
		return nil
	}
	// 日志文件
	logNameArr := strings.Split(logName, ".")
	if len(logNameArr) == 1 {
		logNameArr = append([]string{"logrus"}, logNameArr...)
	}
	fileDir := fmt.Sprintf("./runtime/%s", logNameArr[0])
	filePath := fmt.Sprintf("%s/%s.log", fileDir, strings.Join(logNameArr[1:], "_"))
	// 创建 Logrus 实例
	var logger = logrus.New()
	// 配置 lumberjack 日志轮转，设置日志输出到 lumberjack 管理的文件
	logger.SetOutput(&lumberjack.Logger{
		Filename:   filePath, // 日志文件路径
		MaxSize:    10,       // 单个日志文件最大大小（MB）
		MaxBackups: 3,        // 最多保留的旧日志文件数量
		MaxAge:     28,       // 日志文件最多保留的天数
		Compress:   false,    // 是否压缩旧日志文件
	})
	// 配置输出到标准输出
	// logger.SetOutput(os.Stdout)
	// 设置JSON日志格式
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		}, // 默认字段配置
		TimestampFormat:   time.RFC3339, // 格式化日志时间，注释该行则默认时间格式
		DisableHTMLEscape: true,         // json序列号不编码
		PrettyPrint:       true,         // JSON日志内容会按键值对换行
		// DisableSorting:    true,//禁止字段默认排序
	})
	return logger
}
