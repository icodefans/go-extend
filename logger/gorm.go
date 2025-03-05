package logger

import (
	"bufio"
	"fmt"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// 定义自己的Writer
type MyWriter struct {
	mlog *logrus.Logger
}

// 实现gorm/logger.Writer接口
func (m *MyWriter) Printf(format string, data ...any) {
	// fmt.Println("format:", format)
	myFomart := "%s[%.3fms] [rows:%v] %s"
	logstr := fmt.Sprintf(myFomart, data[0], data[1], data[2], data[3])
	// 利用loggus记录日志
	m.mlog.Info(logstr)
}

func NewMyWriter() *MyWriter {
	logName := "gorm"
	log := logrus.New()
	// 写入空设备，避免控制台输出
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("Open Src File err", err)
	}
	writer := bufio.NewWriter(src)
	log.SetOutput(writer)

	// 日志文件
	fileName := fmt.Sprint("./runtime/", logName, "/", logName, ".log")

	// 设置 rotatelogs
	logWriter, _ := rotatelogs.New(
		// 分割后的文件名称
		fileName+".%Y%m%d.log",
		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),
		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	writeMap := lfshook.WriterMap{
		logrus.TraceLevel: logWriter,
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat:   "2006-01-02 15:04:05", // 格式化日志时间，注释该行则默认时间格式
		DisableHTMLEscape: true,                  // json序列号不编码
	})

	// 新增 Hook
	log.AddHook(lfHook)
	// 是否显示行号
	log.SetReportCaller(false)

	return &MyWriter{mlog: log}
}
