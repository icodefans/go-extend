package redis

import (
	"errors"
	"fmt"

	redisgo "github.com/gomodule/redigo/redis"
	"github.com/redis/go-redis/v9"
)

// 消息队列客户端
type StreamMQClient struct {
	*Redis
	StreamKey    string // stream对应的key值
	GroupName    string // 消费者组名称
	ConsumerName string // 消费者名称
}

// PutMsg 添加消息
func (mqClient *StreamMQClient) PutMsg(streamKey string, maxLen int64, values string) (msgId string, err error) {
	conn, ctx, _ := mqClient.Connect()
	// defer conn.Close()
	// *表示由Redis自己生成消息ID，设置MAXLEN可以保证消息队列的长度不会一直累加
	cmd := conn.XAdd(ctx, &redis.XAddArgs{
		Stream: streamKey,
		MaxLen: maxLen,
		Values: map[string]string{"values": values},
	})
	if msgId, err = cmd.Result(); err != nil {
		return "", fmt.Errorf("XADD failed, err:%s", err)
	}
	return msgId, nil
}

// PutMsgBatch 批量添加消息
func (mqClient *StreamMQClient) PutMsgBatch(streamKey string, maxLen uint32, msgMap map[string]string) (msgId string, err error) {
	if len(msgMap) <= 0 {
		return "", fmt.Errorf("XADD failed, err:%s", "msgMap len <= 0, no need put")
	}

	conn, ctx, _ := mqClient.Connect()
	// defer conn.Close()

	vecMsg := make([]string, 0)
	for msgKey, msgValue := range msgMap {
		vecMsg = append(vecMsg, msgKey)
		vecMsg = append(vecMsg, msgValue)
	}

	cmd := conn.Do(ctx, redisgo.Args{"XADD", streamKey, "MAXLEN", "=", maxLen, "*"}.AddFlat(vecMsg)...)
	if value, err := cmd.Result(); err != nil {
		return "", fmt.Errorf("XADD failed, err:%s", err)
	} else if id, ok := value.(string); !ok {
		return "", fmt.Errorf("数据类型断言失败")
	} else {
		msgId = id
	}
	return msgId, nil
}

func (mqClient *StreamMQClient) ConvertVecInterface(vecReply []any) (msgMap map[string]map[string][]string) {
	msgMap = make(map[string]map[string][]string)
	for keyIndex := 0; keyIndex < len(vecReply); keyIndex++ {
		var keyInfo = vecReply[keyIndex].([]any)
		var key = string(keyInfo[0].([]byte))
		var idList = keyInfo[1].([]any)

		// fmt.Println("StreamKey:", key)
		msgInfoMap := make(map[string][]string)
		for idIndex := 0; idIndex < len(idList); idIndex++ {
			var idInfo = idList[idIndex].([]any)
			var id = string(idInfo[0].([]byte))

			var fieldList = idInfo[1].([]any)
			vecMsg := make([]string, 0)
			for msgIndex := 0; msgIndex < len(fieldList); msgIndex = msgIndex + 2 {
				var msgKey = string(fieldList[msgIndex].([]byte))
				var msgVal = string(fieldList[msgIndex+1].([]byte))
				vecMsg = append(vecMsg, msgKey)
				vecMsg = append(vecMsg, msgVal)
				// fmt.Println("MsgId:", id, "MsgKey:", msgKey, "MsgVal:", msgVal)
			}
			msgInfoMap[id] = vecMsg
		}
		msgMap[key] = msgInfoMap
	}
	return
}

// GetMsgBlock 阻塞方式读取消息
func (mqClient *StreamMQClient) GetMsgBlock(streamKey string, blockSec int32, msgAmount int32) (msgMap map[string]map[string][]string, err error) {
	conn, ctx, _ := mqClient.Connect()
	// defer conn.Close()
	// 在阻塞模式中，可以使用$，表示最新的消息ID（在非阻塞模式下$无意义）
	cmd := conn.Do(ctx, "XREAD", "COUNT", msgAmount, "BLOCK", blockSec*1000, "STREAMS", streamKey, "$")
	var data []any
	if reply, err := cmd.Result(); err != nil {
		return nil, err
	} else if value, ok := reply.([]any); ok {
		return nil, fmt.Errorf("结果数据类型断言失败")
	} else {
		data = value
	}
	// 返回消息转换
	msgMap = mqClient.ConvertVecInterface(data)
	return msgMap, nil
}

// GetMsg 非阻塞方式读取消息
func (mqClient *StreamMQClient) GetMsg(streamKey string, msgAmount int32, beginMsgId string) (
	msgMap map[string]map[string][]string, err error) {
	conn, ctx, _ := mqClient.Connect()
	// defer conn.Close()
	// 从消息ID=beginMsgId往后开始读取，不包含beginMsgId的消息
	cmd := conn.Do(ctx, "XREAD", "COUNT", msgAmount, "STREAMS", streamKey, beginMsgId)
	var data []any
	if reply, err := cmd.Result(); err != nil {
		return nil, err
	} else if value, ok := reply.([]any); !ok {
		return nil, fmt.Errorf("结果数据类型断言失败")
	} else {
		data = value
	}
	// 返回消息转换
	msgMap = mqClient.ConvertVecInterface(data)
	return msgMap, nil
}

// DelMsg 删除消息
func (mqClient *StreamMQClient) DelMsg(streamKey string, msgId string) (err error) {
	if msgId == "" {
		return fmt.Errorf("vecMsgId no need del")
	}
	conn, ctx, _ := mqClient.Connect()
	// defer conn.Close()
	_, err = conn.Do(ctx, "XDEL", streamKey, msgId).Result()
	if err != nil {
		return fmt.Errorf("XDEL failed, msgId:%s err:%s", msgId, err)
	}
	return nil
}

// ReplyAck 返回ACK
func (mqClient *StreamMQClient) ReplyAck(streamKey string, groupName string, vecMsgId []string) (err error) {
	if len(vecMsgId) <= 0 {
		return fmt.Errorf("vecMsgId len <= 0, no need ack")
	}
	conn, ctx, _ := mqClient.Connect()
	// defer conn.Close()
	_, err = conn.Do(ctx, redisgo.Args{"XACK", streamKey, groupName}.AddFlat(vecMsgId)...).Result()
	if err != nil {
		return fmt.Errorf("XACK failed, msgId:%s err:%s", vecMsgId, err)
	}
	return nil
}

// CreateConsumerGroup 创建消费者组
func (mqClient *StreamMQClient) CreateConsumerGroup(streamKey string, groupName string, beginMsgId string) (err error) {
	conn, ctx, _ := mqClient.Connect()
	// defer conn.Close()
	// 最后一个参数表示该组从消息ID=beginMsgId往后开始消费，不包含beginMsgId的消息，如果指定了MKSTREAM， 当stream不存在时，根据key值创建新的STREAM。
	_, err = conn.Do(ctx, "XGROUP", "CREATE", streamKey, groupName, beginMsgId, "MKSTREAM").Result()
	if err != nil {
		return fmt.Errorf("XGROUP CREATE Failed. err:%s", err)
	}
	return nil
}

// DestroyConsumerGroup 销毁消费者组
func (mqClient *StreamMQClient) DestroyConsumerGroup(streamKey string, groupName string) (err error) {
	conn, ctx, _ := mqClient.Connect()
	// defer conn.Close()

	_, err = conn.Do(ctx, "XGROUP", "DESTROY", streamKey, groupName).Result()
	if err != nil {
		return fmt.Errorf("XGROUP DESTROY Failed. err:%s", err)
	}
	return nil
}

// GetMsgByGroupConsumer 组内消息分配操作，组内每个消费者消费多少消息
func (mqClient *StreamMQClient) GetMsgByGroupConsumer(streamKey string, groupName string, consumerName string, count int64) (messages map[string]any, err error) {
	conn, ctx, _ := mqClient.Connect()
	// defer conn.Close()

	// 代表当前消费者还没读取的消息
	var result []redis.XStream
	if result, err = conn.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    groupName,
		Consumer: consumerName,
		Streams:  []string{streamKey, ">"},
		Count:    count,
		Block:    0,
		NoAck:    true,
	}).Result(); err != nil {
		return nil, err
	} else if err != nil && !errors.Is(err, conn.KeyNil) {
		return nil, fmt.Errorf("XREADGROUP failed, err: %s", err)
	} else if len(result) == 0 {
		return nil, nil
	}
	messages = make(map[string]any)
	for _, message := range result[0].Messages {
		if value, ok := message.Values["values"]; ok {
			messages[message.ID] = value
		}
	}
	// 返回消息转换
	return messages, nil
}

// CreateConsumer 创建消费者
func (mqClient *StreamMQClient) CreateConsumer(streamKey string, groupName string, consumerName string) (err error) {
	conn, ctx, _ := mqClient.Connect()
	// defer conn.Close()
	_, err = conn.Do(ctx, "XGROUP", "CREATECONSUMER", streamKey, groupName, consumerName).Result()
	if err != nil {
		return fmt.Errorf("XGROUP CREATECONSUMER Failed. err:%s", err)
	}
	return nil
}

// DelConsumer 删除消费者
func (mqClient *StreamMQClient) DelConsumer(streamKey string, groupName string, consumerName string) (err error) {
	conn, ctx, _ := mqClient.Connect()
	// defer conn.Close()
	_, err = conn.Do(ctx, "XGROUP", "DELCONSUMER", streamKey, groupName, consumerName).Result()
	if err != nil {
		return fmt.Errorf("XGROUP DELCONSUMER Failed. err:%s", err)
	}
	return nil
}

// GetMsgByGroupConsumer 组内消息分配操作，组内每个消费者消费多少消息
func (mqClient *StreamMQClient) GetMsgBlockByGroupConsumer(streamKey string, groupName string, consumerName string, msgAmount, blockSec int32) (msgMap map[string]map[string][]string, err error) {
	conn, ctx, _ := mqClient.Connect()
	// defer conn.Close()

	// 代表当前消费者还没读取的消息
	var data []any
	cmd := conn.Do(ctx, "XREADGROUP", "GROUP", groupName, consumerName, "COUNT", msgAmount, "BLOCK", blockSec*1000, "STREAMS", streamKey, ">")
	if reply, err := cmd.Result(); err != nil && !errors.Is(err, conn.KeyNil) {
		return nil, fmt.Errorf("BLOCK XREADGROUP failed, err: %s", err)
	} else if value, ok := reply.([]any); !ok {
		return nil, fmt.Errorf("结果数据类型断言失败")
	} else {
		data = value
	}
	// 返回消息转换
	msgMap = mqClient.ConvertVecInterface(data)
	return msgMap, nil
}

// GetPendingList 获取等待列表(读取但还未消费的消息)
func (mqClient *StreamMQClient) GetPendingList(streamKey string, groupName string, consumerName string, msgAmount int32) (vecPendingMsg []*PendingMsgInfo, err error) {
	conn, ctx, _ := mqClient.Connect()
	// defer conn.Close()

	var data []any
	cmd := conn.Do(ctx, "XPENDING", streamKey, groupName, "-", "+", msgAmount, consumerName)
	if reply, err := cmd.Result(); err != nil && !errors.Is(err, conn.KeyNil) {
		return nil, fmt.Errorf("XPENDING failed, err: %s", err)
	} else if value, ok := reply.([]any); !ok {
		return nil, fmt.Errorf("结果数据类型断言失败")
	} else {
		data = value
	}
	for iIndex := 0; iIndex < len(data); iIndex++ {
		var msgInfo = data[iIndex].([]any)
		var msgId = string(msgInfo[0].([]byte))
		var belongConsumer = string(msgInfo[1].([]byte))
		var idleTime = msgInfo[2].(int64)
		var readCount = msgInfo[3].(int64)
		pendingMsg := &PendingMsgInfo{msgId, belongConsumer, int(idleTime), int(readCount)}
		vecPendingMsg = append(vecPendingMsg, pendingMsg)
	}
	return vecPendingMsg, nil
}

// MoveMsg 转移消息到其他等待列表中
func (mqClient *StreamMQClient) MoveMsg(streamKey string, groupName string, consumerName string, idleTime int, vecMsgId []string) (err error) {
	if len(vecMsgId) <= 0 {
		return fmt.Errorf("vecMsgId len <= 0, no need move")
	}
	conn, ctx, _ := mqClient.Connect()
	// defer conn.Close()

	_, err = conn.Do(ctx, redisgo.Args{"XCLAIM", streamKey, groupName, consumerName, idleTime * 1000}.AddFlat(vecMsgId)...).Result()
	if err != nil {
		return fmt.Errorf("XCLAIM failed, msgId:%s err:%s", vecMsgId, err)
	}
	return nil
}

// DelDeadMsg 删除不能被消费者处理，也就是不能被 XACK，长时间处于 Pending 列表中的消息
func (mqClient *StreamMQClient) DelDeadMsg(streamKey string, groupName string, vecMsgId []string) (err error) {
	if len(vecMsgId) <= 0 {
		return fmt.Errorf("vecMsgId len <= 0, no need del")
	}
	conn, ctx, _ := mqClient.Connect()
	// defer conn.Close()

	// 删除消息
	_, err = conn.Do(ctx, redisgo.Args{"XDEL", streamKey}.AddFlat(vecMsgId)...).Result()
	if err != nil {
		return fmt.Errorf("XDEL failed, msgId:%s err:%s", vecMsgId, err)
	}
	// 设置ACK，否则消息还会存在pending list中
	_, err = conn.Do(ctx, redisgo.Args{"XACK", streamKey, groupName}.AddFlat(vecMsgId)...).Result()
	if err != nil {
		return fmt.Errorf("XACK failed, groupName:%s msgId:%s err:%s", groupName, vecMsgId, err)
	}
	return nil
}

// GetStreamsLen 获取消息队列的长度，消息消费之后会做标记，不会删除
func (mqClient *StreamMQClient) GetStreamsLen(streamKey string) (number int, err error) {
	conn, ctx, _ := mqClient.Connect()
	// defer conn.Close()

	reply, err := conn.Do(ctx, "XLEN", streamKey).Result()
	if err != nil {
		return 0, fmt.Errorf("XLEN failed, err:%s", err)
	} else if value, ok := reply.(int); !ok {
		return 0, fmt.Errorf("数据类型断言失败")
	} else {
		number = value
	}
	return number, err
}

// MonitorMqInfo 监控服务器队列信息
func (mqClient *StreamMQClient) MonitorMqInfo(streamKey string) (streamMQInfo *StreamMQInfo, err error) {
	conn, ctx, _ := mqClient.Connect()
	// defer conn.Close()

	var data []any
	cmd := conn.Do(ctx, "XINFO", "STREAM", streamKey)
	if reply, err := cmd.Result(); err != nil {
		return nil, fmt.Errorf("XINFO STREAM failed, err:%s", err)
	} else if value, ok := reply.([]any); !ok {
		return nil, fmt.Errorf("结果数据类型断言失败")
	} else if len(value) == 0 {
		return nil, fmt.Errorf("列表记录为空")
	} else {
		data = value
	}

	streamMQInfo = &StreamMQInfo{}
	streamMQInfo.Length = data[1].(int64)
	streamMQInfo.RedixTreeKeys = data[3].(int64)
	streamMQInfo.RedixTreeNodes = data[5].(int64)
	streamMQInfo.LastGeneratedId = string(data[7].([]byte))
	streamMQInfo.Groups, _ = data[9].(int64)

	firstEntryInfo := data[11].([]any)
	firstEntryMsgId := string(firstEntryInfo[0].([]byte))
	vecFirstEntryMsg := firstEntryInfo[1].([]any)
	firstMsgMap := make(map[string]string)
	for iIndex := 0; iIndex < len(vecFirstEntryMsg); iIndex = iIndex + 2 {
		msgKey := string(vecFirstEntryMsg[iIndex].([]byte))
		msgVal := string(vecFirstEntryMsg[iIndex+1].([]byte))
		firstMsgMap[msgKey] = msgVal
	}
	firstEntry := map[string]map[string]string{
		firstEntryMsgId: firstMsgMap,
	}
	streamMQInfo.FirstEntry = &firstEntry

	lastEntryInfo := data[13].([]any)
	lastEntryMsgId := string(lastEntryInfo[0].([]byte))
	vecLastEntryMsg := lastEntryInfo[1].([]any)
	lastMsgMap := make(map[string]string)
	for iIndex := 0; iIndex < len(vecLastEntryMsg); iIndex = iIndex + 2 {
		msgKey := string(vecLastEntryMsg[iIndex].([]byte))
		msgVal := string(vecLastEntryMsg[iIndex+1].([]byte))
		lastMsgMap[msgKey] = msgVal
	}
	lastEntry := map[string]map[string]string{
		lastEntryMsgId: lastMsgMap,
	}
	streamMQInfo.LastEntry = &lastEntry
	return
}

// MonitorConsumerGroupInfo 监控消费者组信息
func (mqClient *StreamMQClient) MonitorConsumerGroupInfo(streamKey string) (groupInfo *GroupInfo, err error) {
	conn, ctx, _ := mqClient.Connect()
	// defer conn.Close()

	var data []any
	reply, err := conn.Do(ctx, "XINFO", "GROUPS", streamKey).Result()
	if err != nil {
		return nil, fmt.Errorf("XINFO STREAM failed, err:%s", err)
	} else if value, ok := reply.([]any); !ok {
		return nil, fmt.Errorf("结果数据类型断言失败")
	} else if len(value) == 0 {
		return nil, fmt.Errorf("列表记录为空")
	} else {
		data = value
	}

	oGroupInfo := data[0].([]any)
	name := string(oGroupInfo[1].([]byte))
	consumers := oGroupInfo[3].(int64)
	pending := oGroupInfo[5].(int64)
	lastDeliveredId := string(oGroupInfo[7].([]byte))
	groupInfo = &GroupInfo{name, consumers, pending, lastDeliveredId}
	return
}

// MonitorConsumerInfo 监控消费者信息
func (mqClient *StreamMQClient) MonitorConsumerInfo(streamKey string, groupName string) (vecConsumerInfo []*ConsumerInfo, err error) {
	conn, ctx, _ := mqClient.Connect()
	// defer conn.Close()

	var data []any
	reply, err := conn.Do(ctx, "XINFO", "CONSUMERS", streamKey, groupName).Result()
	if err != nil {
		return nil, fmt.Errorf("XINFO CONSUMERS failed, err:%s", err)
	} else if value, ok := reply.([]any); !ok {
		return nil, fmt.Errorf("结果数据类型断言失败")
	} else {
		data = value
	}

	for iIndex := 0; iIndex < len(data); iIndex++ {
		oConsumerInfo := data[iIndex].([]any)
		name := string(oConsumerInfo[1].([]byte))
		pending := oConsumerInfo[3].(int64)
		idle := oConsumerInfo[5].(int64)
		vecConsumerInfo = append(vecConsumerInfo, &ConsumerInfo{name, pending, idle})
	}
	return
}
