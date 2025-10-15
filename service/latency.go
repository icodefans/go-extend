package service

import (
	"time"
)

// 程序执行时间
type LatencyTime struct {
	BeginTime time.Time `label:"开始时间"`
	EndTime   time.Time `label:"结束时间"`
}
