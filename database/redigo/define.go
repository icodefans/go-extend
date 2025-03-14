package redigo

const (
	// STREAM_MQ_MAX_LEN  = 500000 // 消息队列最大长度
	STREAM_MQ_MAX_LEN  = 10000 // 消息队列最大长度
	READ_MSG_AMOUNT    = 1000  // 每次读取消息的条数
	READ_MSG_BLOCK_SEC = 30    // 阻塞读取消息时间
	TEST_STREAM_KEY    = "TestStreamKey1"
)
