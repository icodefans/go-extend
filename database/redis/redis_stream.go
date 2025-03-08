package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
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
		NoAck:    noAck,                                   // 默认情况下，XREADGROUP 读取消息后会自动确认消息。使用 NOACK 参数可禁止自动确认，需要后续手动使用 XACK 命令确认消息。
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
func (r *Redis) StreamGroupCreate(streamKey string, groupName string, beginMsgId string) (err error) {
	conn, ctx, _ := r.Connect()
	// 最后一个参数表示该组从消息ID=beginMsgId往后开始消费，不包含beginMsgId的消息，如果指定了MKSTREAM， 当stream不存在时，根据key值创建新的STREAM。
	_, err = conn.Do(ctx, "XGROUP", "CREATE", streamKey, groupName, beginMsgId, "MKSTREAM").Result()
	if err != nil {
		return fmt.Errorf("XGROUP CREATE Failed. err:%s", err)
	}
	return nil
}

// 组消费者创建
func (r *Redis) StreamGroupConsumerCreate(streamKey string, groupName string, consumerName string) (err error) {
	conn, ctx, _ := r.Connect()
	_, err = conn.Do(ctx, "XGROUP", "CREATECONSUMER", streamKey, groupName, consumerName).Result()
	if err != nil {
		return fmt.Errorf("XGROUP CREATECONSUMER Failed. err:%s", err)
	}
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
