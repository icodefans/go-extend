package redigo

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// 消息队列客户端
type StreamMQClient struct {
	ConnPool     *redis.Pool
	StreamKey    string // stream对应的key值
	GroupName    string // 消费者组名称
	ConsumerName string // 消费者名称
}

// PutMsg 添加消息
func (mqClient *StreamMQClient) PutMsg(streamKey string, maxLen uint32, msgKey string, msgValue string) (strMsgId string, err error) {
	conn := mqClient.ConnPool.Get()
	defer conn.Close()
	// *表示由Redis自己生成消息ID，设置MAXLEN可以保证消息队列的长度不会一直累加
	strMsgId, err = redis.String(conn.Do("XADD",
		streamKey, "MAXLEN", "=", maxLen, "*", msgKey, msgValue))
	if err != nil {
		fmt.Println("XADD failed, err: ", err)
		return "", err
	}
	// fmt.Println("Reply Msg Id:", strMsgId)
	return strMsgId, nil
}

// PutMsgBatch 批量添加消息
func (mqClient *StreamMQClient) PutMsgBatch(streamKey string, maxLen uint32, msgMap map[string]string) (msgId string, err error) {
	if len(msgMap) <= 0 {
		fmt.Println("msgMap len <= 0, no need put")
		return msgId, nil
	}

	conn := mqClient.ConnPool.Get()
	defer conn.Close()

	vecMsg := make([]string, 0)
	for msgKey, msgValue := range msgMap {
		vecMsg = append(vecMsg, msgKey)
		vecMsg = append(vecMsg, msgValue)
	}

	msgId, err = redis.String(conn.Do("XADD",
		redis.Args{streamKey, "MAXLEN", "=", maxLen, "*"}.AddFlat(vecMsg)...))
	if err != nil {
		fmt.Println("XADD failed, err: ", err)
		return "", err
	}

	// fmt.Println("Reply Msg Id:", msgId)
	return msgId, nil
}

func (mqClient *StreamMQClient) ConvertVecInterface(vecReply []any) (msgMap map[string]map[string][]string) {
	msgMap = make(map[string]map[string][]string, 0)
	for keyIndex := 0; keyIndex < len(vecReply); keyIndex++ {
		var keyInfo = vecReply[keyIndex].([]any)
		var key = string(keyInfo[0].([]byte))
		var idList = keyInfo[1].([]any)

		// fmt.Println("StreamKey:", key)
		msgInfoMap := make(map[string][]string, 0)
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
func (mqClient *StreamMQClient) GetMsgBlock(streamKey string, blockSec int32, msgAmount int32) (
	msgMap map[string]map[string][]string, err error) {
	conn := mqClient.ConnPool.Get()
	defer conn.Close()
	// 在阻塞模式中，可以使用$，表示最新的消息ID（在非阻塞模式下$无意义）
	reply, err := redis.Values(conn.Do("XREAD",
		"COUNT", msgAmount, "BLOCK", blockSec*1000, "STREAMS", streamKey, "$"))
	if err != nil && err != redis.ErrNil {
		fmt.Println("BLOCK XREAD failed, err: ", err)
		return nil, err
	}

	// 返回消息转换
	msgMap = mqClient.ConvertVecInterface(reply)
	fmt.Println("MsgMap:", msgMap)
	return msgMap, nil
}

// GetMsg 非阻塞方式读取消息
func (mqClient *StreamMQClient) GetMsg(streamKey string, msgAmount int32, beginMsgId string) (
	msgMap map[string]map[string][]string, err error) {
	conn := mqClient.ConnPool.Get()
	defer conn.Close()
	// 从消息ID=beginMsgId往后开始读取，不包含beginMsgId的消息
	reply, err := redis.Values(conn.Do("XREAD", "COUNT", msgAmount, "STREAMS", streamKey, beginMsgId))
	if err != nil && err != redis.ErrNil {
		fmt.Println("XREAD failed, err: ", err)
		return nil, err
	}

	// 返回消息转换
	msgMap = mqClient.ConvertVecInterface(reply)
	fmt.Println("MsgMap:", msgMap)
	return msgMap, nil
}

// DelMsg 删除消息
func (mqClient *StreamMQClient) DelMsg(streamKey string, vecMsgId []string) (err error) {
	if len(vecMsgId) <= 0 {
		fmt.Println("vecMsgId len <= 0, no need del")
		return nil
	}

	conn := mqClient.ConnPool.Get()
	defer conn.Close()

	for _, msgId := range vecMsgId {
		_, err := redis.Int(conn.Do("XDEL", streamKey, msgId))
		if err != nil {
			fmt.Println("XDEL failed, msgId:", msgId, "err:", err)
		}
	}
	return nil
}

// ReplyAck 返回ACK
func (mqClient *StreamMQClient) ReplyAck(streamKey string, groupName string, vecMsgId []string) error {
	if len(vecMsgId) <= 0 {
		// fmt.Println("vecMsgId len <= 0, no need ack")
		return nil
	}

	conn := mqClient.ConnPool.Get()
	defer conn.Close()

	// fmt.Println("Start ReplyAck, vecMsgId:", vecMsgId)
	_, err := redis.Int(conn.Do("XACK", redis.Args{streamKey, groupName}.AddFlat(vecMsgId)...))
	if err != nil {
		fmt.Println("XACK failed, msgId:", vecMsgId, "err:", err)
		return err
	}
	// fmt.Println("ReplyAck Success")
	return nil
}

// CreateConsumerGroup 创建消费者组
func (mqClient *StreamMQClient) CreateConsumerGroup(streamKey string, groupName string, beginMsgId string) error {
	conn := mqClient.ConnPool.Get()
	defer conn.Close()
	// 最后一个参数表示该组从消息ID=beginMsgId往后开始消费，不包含beginMsgId的消息，如果指定了MKSTREAM， 当stream不存在时，根据key值创建新的STREAM。
	_, err := redis.String(conn.Do("XGROUP", "CREATE", streamKey, groupName, beginMsgId, "MKSTREAM"))
	if err != nil {
		// fmt.Println("XGROUP CREATE Failed. err:", err)
		return err
	}
	return nil
}

// DestroyConsumerGroup 销毁消费者组
func (mqClient *StreamMQClient) DestroyConsumerGroup(streamKey string, groupName string) error {
	conn := mqClient.ConnPool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("XGROUP", "DESTROY", streamKey, groupName))
	if err != nil {
		fmt.Println("XGROUP DESTROY Failed. err:", err)
		return err
	}
	return nil
}

// GetMsgByGroupConsumer 组内消息分配操作，组内每个消费者消费多少消息
func (mqClient *StreamMQClient) GetMsgByGroupConsumer(streamKey string, groupName string, consumerName string, msgAmount int32) (msgMap map[string]map[string][]string, err error) {
	conn := mqClient.ConnPool.Get()
	defer conn.Close()

	// >代表当前消费者还没读取的消息
	reply, err := redis.Values(conn.Do("XREADGROUP",
		"GROUP", groupName, consumerName, "COUNT", msgAmount, "STREAMS", streamKey, ">"))
	if err != nil && err != redis.ErrNil {
		fmt.Println("XREADGROUP failed, err: ", err)
		return nil, err
	}

	// 返回消息转换
	msgMap = mqClient.ConvertVecInterface(reply)
	// fmt.Println("MsgMap:", msgMap)
	return msgMap, nil
}

// CreateConsumer 创建消费者
func (mqClient *StreamMQClient) CreateConsumer(streamKey string, groupName string, consumerName string) error {
	conn := mqClient.ConnPool.Get()
	defer conn.Close()

	_, err := redis.Int64(conn.Do("XGROUP", "CREATECONSUMER", streamKey, groupName, consumerName))
	if err != nil {
		// fmt.Println("XGROUP CREATECONSUMER Failed. err:", err)
		return err
	}
	return nil
}

// DelConsumer 删除消费者
func (mqClient *StreamMQClient) DelConsumer(streamKey string, groupName string, consumerName string) error {
	conn := mqClient.ConnPool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("XGROUP", "DELCONSUMER", streamKey, groupName, consumerName))
	if err != nil {
		fmt.Println("XGROUP DELCONSUMER Failed. err:", err)
		return err
	}
	return nil
}

// GetMsgByGroupConsumer 组内消息分配操作，组内每个消费者消费多少消息
func (mqClient *StreamMQClient) GetMsgBlockByGroupConsumer(streamKey string, groupName string,
	consumerName string, msgAmount, blockSec int32) (msgMap map[string]map[string][]string, err error) {
	conn := mqClient.ConnPool.Get()
	defer conn.Close()

	// >代表当前消费者还没读取的消息
	reply, err := redis.Values(conn.Do("XREADGROUP", "GROUP", groupName,
		consumerName, "COUNT", msgAmount, "BLOCK", blockSec*1000, "STREAMS", streamKey, ">"))
	if err != nil && err != redis.ErrNil {
		fmt.Println("BLOCK XREADGROUP failed, err: ", err)
		return nil, err
	}

	// 返回消息转换
	msgMap = mqClient.ConvertVecInterface(reply)
	fmt.Println("MsgMap:", msgMap)
	return msgMap, nil
}

// GetPendingList 获取等待列表(读取但还未消费的消息)
func (mqClient *StreamMQClient) GetPendingList(streamKey string, groupName string, consumerName string, msgAmount int32) (
	vecPendingMsg []*PendingMsgInfo, err error) {
	conn := mqClient.ConnPool.Get()
	defer conn.Close()

	reply, err := redis.Values(conn.Do("XPENDING", streamKey, groupName, "-", "+", msgAmount, consumerName))
	if err != nil {
		fmt.Println("XPENDING failed, err: ", err)
		return nil, err
	}

	for iIndex := 0; iIndex < len(reply); iIndex++ {

		var msgInfo = reply[iIndex].([]any)
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
func (mqClient *StreamMQClient) MoveMsg(streamKey string, groupName string, consumerName string, idleTime int, vecMsgId []string) error {
	if len(vecMsgId) <= 0 {
		fmt.Println("vecMsgId len <= 0, no need move")
		return nil
	}

	conn := mqClient.ConnPool.Get()
	defer conn.Close()

	_, err := redis.Values(conn.Do("XCLAIM", redis.Args{streamKey, groupName, consumerName, idleTime * 1000}.AddFlat(vecMsgId)...))
	if err != nil {
		fmt.Println("XCLAIM failed, msgId:", vecMsgId, "err:", err)
		return err
	}
	return nil
}

// DelDeadMsg 删除不能被消费者处理，也就是不能被 XACK，长时间处于 Pending 列表中的消息
func (mqClient *StreamMQClient) DelDeadMsg(streamKey string, groupName string, vecMsgId []string) error {
	if len(vecMsgId) <= 0 {
		fmt.Println("vecMsgId len <= 0, no need del")
		return nil
	}

	conn := mqClient.ConnPool.Get()
	defer conn.Close()

	// 删除消息
	_, err := redis.Int(conn.Do("XDEL", redis.Args{streamKey}.AddFlat(vecMsgId)...))
	if err != nil {
		fmt.Println("XDEL failed, msgId:", vecMsgId, "err:", err)
		return err
	}
	// 设置ACK，否则消息还会存在pending list中
	_, err = redis.Int(conn.Do("XACK", redis.Args{streamKey, groupName}.AddFlat(vecMsgId)...))
	if err != nil {
		fmt.Println("XACK failed, groupName:", groupName, "msgId:", vecMsgId, "err:", err)
		return err
	}
	return nil
}

// GetStreamsLen 获取消息队列的长度，消息消费之后会做标记，不会删除
func (mqClient *StreamMQClient) GetStreamsLen(streamKey string) int {
	conn := mqClient.ConnPool.Get()
	defer conn.Close()

	reply, err := redis.Int(conn.Do("XLEN", streamKey))
	if err != nil {
		fmt.Println("XLEN failed, err:", err)
		return -1
	}
	return reply
}

// MonitorMqInfo 监控服务器队列信息
func (mqClient *StreamMQClient) MonitorMqInfo(streamKey string) (streamMQInfo *StreamMQInfo) {
	conn := mqClient.ConnPool.Get()
	defer conn.Close()

	reply, err := redis.Values(conn.Do("XINFO", "STREAM", streamKey))
	if err != nil || len(reply) <= 0 {
		fmt.Println("XINFO STREAM failed, err:", err)
		return nil
	}
	fmt.Println("reply len:", len(reply))

	streamMQInfo = &StreamMQInfo{}
	streamMQInfo.Length = reply[1].(int64)
	streamMQInfo.RedixTreeKeys = reply[3].(int64)
	streamMQInfo.RedixTreeNodes = reply[5].(int64)
	streamMQInfo.LastGeneratedId = string(reply[7].([]byte))
	streamMQInfo.Groups, _ = reply[9].(int64)

	firstEntryInfo := reply[11].([]any)
	firstEntryMsgId := string(firstEntryInfo[0].([]byte))
	vecFirstEntryMsg := firstEntryInfo[1].([]any)
	firstMsgMap := make(map[string]string, 0)
	for iIndex := 0; iIndex < len(vecFirstEntryMsg); iIndex = iIndex + 2 {
		msgKey := string(vecFirstEntryMsg[iIndex].([]byte))
		msgVal := string(vecFirstEntryMsg[iIndex+1].([]byte))
		firstMsgMap[msgKey] = msgVal
	}
	firstEntry := map[string]map[string]string{
		firstEntryMsgId: firstMsgMap,
	}
	streamMQInfo.FirstEntry = &firstEntry

	lastEntryInfo := reply[13].([]any)
	lastEntryMsgId := string(lastEntryInfo[0].([]byte))
	vecLastEntryMsg := lastEntryInfo[1].([]any)
	lastMsgMap := make(map[string]string, 0)
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
func (mqClient *StreamMQClient) MonitorConsumerGroupInfo(streamKey string) (groupInfo *GroupInfo) {
	conn := mqClient.ConnPool.Get()
	defer conn.Close()

	reply, err := redis.Values(conn.Do("XINFO", "GROUPS", streamKey))
	if err != nil || len(reply) <= 0 {
		fmt.Println("XINFO GROUPS failed, err:", err)
		return nil
	}
	fmt.Println("reply len:", len(reply))

	oGroupInfo := reply[0].([]any)
	name := string(oGroupInfo[1].([]byte))
	consumers := oGroupInfo[3].(int64)
	pending := oGroupInfo[5].(int64)
	lastDeliveredId := string(oGroupInfo[7].([]byte))
	groupInfo = &GroupInfo{name, consumers, pending, lastDeliveredId}

	return
}

// MonitorConsumerInfo 监控消费者信息
func (mqClient *StreamMQClient) MonitorConsumerInfo(streamKey string, groupName string) (vecConsumerInfo []*ConsumerInfo) {
	conn := mqClient.ConnPool.Get()
	defer conn.Close()

	reply, err := redis.Values(conn.Do("XINFO", "CONSUMERS", streamKey, groupName))
	if err != nil {
		fmt.Println("XINFO CONSUMERS failed, err:", err)
		return nil
	}
	fmt.Println("reply len:", len(reply))

	for iIndex := 0; iIndex < len(reply); iIndex++ {
		oConsumerInfo := reply[iIndex].([]any)
		name := string(oConsumerInfo[1].([]byte))
		pending := oConsumerInfo[3].(int64)
		idle := oConsumerInfo[5].(int64)
		vecConsumerInfo = append(vecConsumerInfo, &ConsumerInfo{name, pending, idle})
	}
	return
}
