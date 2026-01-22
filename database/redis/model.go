package redis

import (
	"time"
)

// 队列消费组消费者清理配置
type StreamCleanupConfig struct {
	ConsumerMaxIdleTime time.Duration // 消费者最大空闲时间
	ConsumerMinPending  int64         // 消费者最小待处理消息阈值
	GroupMaxIdleTime    time.Duration // 消费组最大空闲时间
	GroupMinPending     int64         // 消费组最小待处理消息阈值
	GroupEmptyDelete    bool          // 消费组为空是否删除
}
