// Redis
package redis

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
	KeyNil   error  // 定义键不存在错误，避免业务中导入redis包
	*redis.Client
}

// 实例化连接
func (r *Redis) Connect() (rdb *Redis, ctx context.Context, cancel context.CancelFunc) {
	// 参数设置
	r.KeyNil = redis.Nil
	// 超时控制
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	// 单利模式
	if r.Client != nil {
		return r, ctx, cancel
	}
	// 客户端实例化
	r.Client = redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%d", r.Host, r.Port), // 连接地址
		Username:    "default",                            // 用户名
		Password:    r.Password,                           // 密码
		DB:          int(r.Index),                         // 数据库
		PoolSize:    int(r.PoolSize),                      // 连接池大小
		DialTimeout: 30 * time.Second,                     // 拨号超时
		ReadTimeout: 60 * time.Second,                     // 读取超时
		MaxRetries:  2,                                    // 最大重试次数
	})
	return r, ctx, cancel
}

// 扩展Incr方法，支持过期时间
func (r *Redis) IncrX(ctx context.Context, key string, expire uint32) (val uint32, err error) {
	var incrV int64
	incrCmd := r.Incr(ctx, key)
	if err = incrCmd.Err(); err != nil {
		return 0, err
	} else if incrV, err = incrCmd.Result(); err != nil {
		return 0, err
	} else if incrV == 1 {
		r.Expire(ctx, key, time.Second*time.Duration(expire))
	}
	return uint32(incrV), nil
}
