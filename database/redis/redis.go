// Redis
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis配置
type Redis struct {
	Host     string           `mapstructure:"Host"`     // 连接地址
	Port     int32            `mapstructure:"Port"`     // 连接地址
	Password string           `mapstructure:"Password"` // 密码
	Index    int32            `mapstructure:"Index"`    // 数据库
	PoolSize int32            `mapstructure:"PoolSize"` // 连接池大小
	Timeout  int32            `mapstructure:"Timeout"`  // 超时时间
	KeyNil   error            // 定义键不存在错误，避免业务中导入redis包
	Ctx      *context.Context // 上下文
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
	ctx2 := context.Background()
	r.Ctx = &ctx2
	return r, ctx, cancel
}

// 扩展Incr方法，支持过期时间
func (r *Redis) IncrX(ctx context.Context, key string, expire uint32) (incrV uint32, err error) {
	// Lua script to increment a counter and check if it exceeds the limit
	luaScript := `
        local key = KEYS[1]
        local expiration = tonumber(ARGV[1])
        local counter = redis.call("incr", key)
        if counter == 1 then
            redis.call("expire", key, expiration)
        end
        return counter`
	keys := []string{key}
	args := []string{fmt.Sprintf("%d", expire)}
	res, err := r.Eval(ctx, luaScript, keys, args).Result()
	if err != nil {
		return 0, err
	}
	return uint32(res.(int64)), nil
}

// 扩展Incr方法，支持过期时间
func (r *Redis) IncrY(ctx context.Context, key string, expire uint32) (val uint32, err error) {
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

// 模糊匹配删除
// pattern 模糊匹配规则 "your_pattern*"
func (r *Redis) DelX(ctx context.Context, pattern string) (rows uint32, err error) {
	// 定义 Lua 脚本
	luaScript := `
		local pattern = ARGV[1]
		local cursor = '0'
		local totalDeleted = 0
		repeat
			local result = redis.call('SCAN', cursor, 'MATCH', pattern)
			cursor = result[1]
			local keys = result[2]
			if #keys > 0 then
				local deleted = redis.call('DEL', unpack(keys))
				totalDeleted = totalDeleted + deleted
			end
		until cursor == '0'
		return totalDeleted
	`
	// 执行 Lua 脚本
	res, err := r.Eval(ctx, luaScript, []string{}, pattern).Result()
	if err != nil && err != redis.Nil {
		return 0, fmt.Errorf("执行 Lua 脚本时出错: %v", err)
	}
	return uint32(res.(int64)), nil
}
