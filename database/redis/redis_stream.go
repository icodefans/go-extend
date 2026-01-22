package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// 消息添加
func (r *Redis) StreamMessageAdd(streamKey string, maxLen int64, values string) (msgId string, err error) {
	conn, ctx, _ := r.Connect()
	// *表示由Redis自己生成消息ID，设置MAXLEN可以保证消息队列的长度不会一直累加
	if msgId, err = conn.XAdd(ctx, &redis.XAddArgs{
		Stream: streamKey,
		MaxLen: maxLen,
		Values: map[string]string{"values": values},
	}).Result(); err != nil {
		return "", fmt.Errorf("XADD failed, err:%s", err)
	}
	return msgId, nil
}

// 组内消息分配操作，组内每个消费者消费多少消息
func (r *Redis) StreamMessageByGroupConsumer(streamKey string, groupName string, consumerName string, count int64, block int64, noAck bool) (messages map[string]string, err error) {
	conn, ctx, _ := r.Connect()
	var result []redis.XStream
	if result, err = conn.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    groupName,                               // 消费者组的名称，用于标识一组协同工作的消费者。
		Consumer: consumerName,                            // 当前执行读取操作的消费者名称，在消费者组中唯一标识该消费者。
		Streams:  []string{streamKey, ">"},                // <key>：流的名称，可以指定多个流，以实现从多个流中同时读取消息。<id>：每个流对应的起始消息 ID，通常使用 > 表示从流中未被消费者组处理过的最新消息开始读取。
		Count:    count,                                   // 指定每次读取的最大消息数量。若不指定，Redis 会尝试读取尽可能多的消息。
		Block:    time.Millisecond * time.Duration(block), // 若设置为 0，则表示无限期阻塞，直到有新消息到来；若在指定时间内没有新消息，命令将返回空结果
		NoAck:    noAck,                                   // 读取消息后，会自动将其放入 Pending（等待）列表，等待消费者显式地发送 XACK 命令进行确认
	}).Result(); err != nil && !errors.Is(err, conn.KeyNil) {
		return nil, fmt.Errorf("XREADGROUP failed, err: %s", err)
	} else if len(result) == 0 {
		return nil, nil
	}
	messages = make(map[string]string)
	for _, message := range result[0].Messages {
		if value, ok := message.Values["values"]; !ok {
			continue
		} else if val, ok := value.(string); !ok {
			continue
		} else {
			messages[message.ID] = val
		}
	}
	// 返回消息转换
	return messages, nil
}

// 消费者组创建
// XGROUP CREATE mystream mygroup $ IDLE 120000 PEL-EXPIRY 300000 MKSTREAM
// IDLE 120000： 消费者如果 2 分钟不活动，则其持有的所有 PEL 消息会被重新分配。
// PEL-EXPIRY 300000： 任何消息在 PEL 中停留超过 5 分钟，无论消费者状态如何，都会被自动清除并重新投递。
func (r *Redis) StreamGroupCreate(streamKey string, groupName string, beginMsgId string) (err error) {
	var (
		ctxKey   = fmt.Sprintf("%s:%s", streamKey, groupName)
		ctxValue = false
	)
	if ctxValue, _ = (*r.Ctx).Value(ctxKey).(bool); ctxValue {
		return nil
	}
	// 最后一个参数表示该组从消息ID=beginMsgId往后开始消费，不包含beginMsgId的消息，如果指定了MKSTREAM， 当stream不存在时，根据key值创建新的STREAM。
	conn, ctx, _ := r.Connect()
	if _, err = conn.Do(ctx, "XGROUP", "CREATE", streamKey, groupName, beginMsgId, "MKSTREAM").Result(); err != nil {
		return fmt.Errorf("XGROUP CREATE Failed. err:%s", err)
	}
	*r.Ctx = context.WithValue(*r.Ctx, ctxKey, true)
	return nil
}

// 消费者组销毁
func (r *Redis) StreamGroupDestroy(streamKey string, groupName string) (err error) {
	conn, ctx, _ := r.Connect()
	if _, err = conn.XGroupDestroy(ctx, streamKey, groupName).Result(); err != nil {
		return fmt.Errorf("XGROUP Destroy Failed. err:%s", err)
	}
	return nil
}

// 组消费者创建
func (r *Redis) StreamGroupConsumerCreate(streamKey string, groupName string, consumerName string) (err error) {
	var (
		ctxKey   = fmt.Sprintf("%s:%s", streamKey, groupName)
		ctxValue = false
	)
	if ctxValue, _ = (*r.Ctx).Value(ctxKey).(bool); ctxValue {
		return nil
	}
	conn, ctx, _ := r.Connect()
	if _, err = conn.Do(ctx, "XGROUP", "CREATECONSUMER", streamKey, groupName, consumerName).Result(); err != nil {
		return fmt.Errorf("XGROUP CREATECONSUMER Failed. err:%s", err)
	}
	*r.Ctx = context.WithValue(*r.Ctx, ctxKey, true)
	return nil
}

// 消息已消费ACK确认
func (r *Redis) StreamXAck(streamKey string, groupName string, vecMsgId []string) (err error) {
	if len(vecMsgId) <= 0 {
		return fmt.Errorf("vecMsgId len <= 0, no need ack")
	}
	conn, ctx, _ := r.Connect()
	if _, err = conn.XAck(ctx, streamKey, groupName, vecMsgId...).Result(); err != nil {
		return err
	}
	return nil
}

// 队列消费组消费者清理
func (r *Redis) StreamCleanup(stream string, config StreamCleanupConfig) error {
	groups, err := r.XInfoGroups(*r.Ctx, stream).Result()
	if err != nil {
		return fmt.Errorf("r.XInfoGroups Err:%s", err)
	}

	for _, group := range groups {
		consumers, err := r.XInfoConsumers(*r.Ctx, stream, group.Name).Result()
		if err != nil {
			return fmt.Errorf("r.XInfoConsumers Err:%s", err)
		}

		activeConsumers := 0     // 活跃的消费者统计
		totalPending := int64(0) // Pending记录统计

		// 清理不活跃消费者
		for _, consumer := range consumers {
			idleTime := consumer.Idle // 消费者空闲时间
			totalPending += consumer.Pending
			if idleTime > config.ConsumerMaxIdleTime && consumer.Pending <= config.ConsumerMinPending {
				r.XGroupDelConsumer(*r.Ctx, stream, group.Name, consumer.Name)
				fmt.Printf("删除不活跃消费者: %s (空闲: %v, 待处理: %d)\n", consumer.Name, idleTime, consumer.Pending)
			} else {
				activeConsumers++
			}
		}

		// 检查是否删除空消费组
		if config.GroupEmptyDelete == false || activeConsumers > 0 {
			// next
		} else if group.Pending <= config.GroupMinPending {
			r.XGroupDestroy(*r.Ctx, stream, group.Name)
			fmt.Printf("删除不活跃消费组: %s\n", group.Name)
		}
	}

	return nil
}
