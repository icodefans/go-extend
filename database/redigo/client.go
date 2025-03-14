package redigo

import (
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

// 客户端连接配置
type RedisConnOpt struct {
	Host      string `mapstructure:"Host"`
	Port      int32  `mapstructure:"Port"`
	Password  string `mapstructure:"Password"`
	Index     int32  `mapstructure:"Index"`
	Timeout   int32  `mapstructure:"Timeout"`  // 最大建立连接等待时间
	MaxActive int32  `mapstructure:"PoolSize"` // 连接池的最大数据库连接数。设为0表示无限制。
	MaxIdle   int32  `mapstructure:"MaxIdle"`  // 最大空闲数，数据库连接的最大空闲时间。超过空闲时间，数据库连 接将被标记为不可用，然后被释放。设为0表示无限制。
	MaxWait   int32  `mapstructure:"MaxWait"`  // 最大建立连接等待时间。如果超过此时间将接到异常。设为-1表示 无限制。
	client    *redis.Pool
}

// 消息队列客户端实例化
func (opt *RedisConnOpt) Connect() (client *redis.Pool) {
	if opt.client != client {
		println("aaaa")
		return opt.client
	}
	opt.client = &redis.Pool{
		MaxIdle:     int(opt.MaxIdle), // 3
		IdleTimeout: 240 * time.Second,
		MaxActive:   int(opt.MaxActive), // 10
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", opt.Host, opt.Port))
			if err != nil {
				log.Fatalf("Redis.Dial: %v", err)
				return nil, err
			}
			/*
				if _, err := c.Do("AUTH", opt.Password); err != nil {
					c.Close()
					log.Fatalf("Redis.AUTH: %v", err)
					return nil, err
				}
			*/
			if _, err := c.Do("SELECT", opt.Index); err != nil {
				c.Close()
				log.Fatalf("Redis.SELECT: %v", err)
				return nil, err
			}
			return c, nil
		},
	}
	return opt.client
}

// 消息队列客户端实例化
func NewStreamMQClient(opt RedisConnOpt) *StreamMQClient {
	return &StreamMQClient{
		ConnPool: &redis.Pool{
			MaxIdle:     int(opt.MaxIdle), // 3
			IdleTimeout: 240 * time.Second,
			MaxActive:   int(opt.MaxActive), // 10
			Wait:        true,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", opt.Host, opt.Port))
				if err != nil {
					log.Fatalf("Redis.Dial: %v", err)
					return nil, err
				}
				/*
					if _, err := c.Do("AUTH", opt.Password); err != nil {
						c.Close()
						log.Fatalf("Redis.AUTH: %v", err)
						return nil, err
					}
				*/
				if _, err := c.Do("SELECT", opt.Index); err != nil {
					c.Close()
					log.Fatalf("Redis.SELECT: %v", err)
					return nil, err
				}
				return c, nil
			},
		},
	}
}
