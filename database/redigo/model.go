package redigo

// 等待列表中的消息属性
type PendingMsgInfo struct {
	MsgId          string // 消息ID
	BelongConsumer string // 所属消费者
	IdleTime       int    // 已读取未消费时长
	ReadCount      int    // 消息被读取次数
}

// 消息队列信息
type StreamMQInfo struct {
	Length          int64                         // 消息队列长度
	RedixTreeKeys   int64                         // 基数树key数
	RedixTreeNodes  int64                         // 基数树节点数
	LastGeneratedId string                        // 最后一个生成的消息ID
	Groups          int64                         // 消费组个数
	FirstEntry      *map[string]map[string]string // 第一个消息体
	LastEntry       *map[string]map[string]string // 最后一个消息体
}

// 消费组信息
type GroupInfo struct {
	Name            string // 消费组名称
	Consumers       int64  // 组内消费者个数
	Pending         int64  // 组内所有消费者的等待列表总长度
	LastDeliveredId string // 组内最后一条被消费的消息ID
}

// 消费者信息
type ConsumerInfo struct {
	Name    string // 消费者名称
	Pending int64  // 等待队列长度
	Idle    int64  // 消费者空闲时间（毫秒）
}
