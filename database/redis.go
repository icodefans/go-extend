// Redis
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Redis配置
type Redis struct {
	Host     string `mapstructure:"Host"`     // 连接地址
	Port     int32  `mapstructure:"Port"`     // 连接地址
	Password string `mapstructure:"Password"` // 密码
	Index    int32  `mapstructure:"Index"`    // 数据库
	PoolSize int32  `mapstructure:"PoolSize"` // 连接池大小
	Timeout  int32  `mapstructure:"Timeout"`  // 超时时间
	rdb      *redis.Client
}

// 实例化连接
func (config *Redis) Connect() (rdb *redis.Client, ctx context.Context, cancel context.CancelFunc) {
	// 单利模式
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	if config.rdb != rdb {
		return config.rdb, ctx, cancel
	}
	// 客户端实例化
	config.rdb = redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%d", config.Host, config.Port), // 连接地址
		Username:    "default",                                      // 用户名
		Password:    config.Password,                                // 密码
		DB:          int(config.Index),                              // 数据库
		PoolSize:    int(config.PoolSize),                           // 连接池大小
		DialTimeout: 30 * time.Second,                               // 拨号超时
		ReadTimeout: 60 * time.Second,                               // 读取超时
		MaxRetries:  2,                                              // 最大重试次数
	})
	// 成功返回
	return config.rdb, ctx, cancel
}
