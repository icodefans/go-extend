// 自动缓存
package cache

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/icodefans/go-extend/database/redis"
	"github.com/syyongx/php2go"
)

// 自动缓存调用
func AutoCall(mode Mode, f DataHandlerFunc, data any, args ...any) (err error) {
	// 异常处理
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		rec := recover()
		if rec != nil {
			err = fmt.Errorf("自动缓存调用异常")
			switch rec.(type) {
			case runtime.Error: // 运行时错误
				fmt.Println("runtime error:", rec)
			default: // 非运行时错误
				fmt.Println("error:", rec)
			}
		}
	}()
	// 初始化缓存参数
	init := f(INIT, args...)
	if init.Error != nil {
		return init.Error
	} else if init.Key == "" {
		return errors.New("缓存方法未设置缓存标识")
	} else if init.Redis.Host == "" {
		return errors.New("缓存方法未设置缓存配置")
	} else if init.Data != nil {
		return errors.New("缓存配置未返回")
	}
	var sub_key string
	if len(args) > 0 {
		sub_args, err := json.Marshal(args)
		if err != nil {
			return errors.New("缓存方法参数序列化错误")
		}
		sub_key = fmt.Sprintf(":%x", md5.Sum(sub_args))
	}
	key := fmt.Sprint("auto://", init.Key, sub_key)
	rdb, ctx, _ := init.Redis.Connect()
	// 缓存删除
	if mode == DELETE {
		rdb.Del(ctx, key)
		return nil
	}
	// 缓存手动设置
	if mode == SET && data == nil {
		return errors.New("缓存数据不能为空")
	} else if mode == SET {
		if json_data, err := json.Marshal(data); err != nil {
			return err
		} else if err := autoSet(json_data, key, init.Expire, init.Redis); err != nil {
			return fmt.Errorf("AutoCache Set Error:%s\n", err)
		}
		return nil
	}
	// 缓存获取,支持缓存关闭
	value, err := autoGet(key, init.Redis)
	if err == nil && init.Expire > 0 {
		return json.Unmarshal([]byte(value), data)
	} else if err != nil && !errors.Is(err, rdb.KeyNil) {
		return fmt.Errorf("AutoCache Get Error:%s\n", err)
	}
	// 缓存不存在则调用方法获取
	result := f(READ, args...)
	if result.Error != nil {
		return result.Error
	} else if result.Data == nil {
		return errors.New("缓存数据不能为空")
	}
	// 缓存自动设置
	if json_data, err := json.Marshal(result.Data); err != nil {
		return err
	} else if err = autoSet(json_data, key, result.Expire, result.Redis); err != nil {
		return fmt.Errorf("AutoCache Set Error:%s\n", err)
	} else if err = json.Unmarshal(json_data, data); err != nil {
		return err
	}
	return nil
}

// 自动缓存获取
func autoGet(key string, redis_conf redis.Redis) (json_data string, err error) {
	rdb, ctx, _ := redis_conf.Connect()
	json_data, err = rdb.Get(ctx, key).Result()
	// 数据读取报错
	if err != nil && !errors.Is(err, rdb.KeyNil) {
		return "", err
	}
	// 数据为空处理逻辑
	if err != nil && errors.Is(err, rdb.KeyNil) {
		return "", rdb.KeyNil
	}
	// 缓存不为空，解析缓存
	var data struct {
		TimeOut  int64
		JsonData string
	}
	if err := json.Unmarshal([]byte(json_data), &data); err != nil {
		return "", err
	}
	// 缓存未过期，直接返回
	if time.Now().Unix() <= data.TimeOut {
		return data.JsonData, nil
	}
	// 成功返回
	return "", rdb.KeyNil
}

// 自动缓存设置
func autoSet(json_data []byte, key string, expire uint32, redis_conf redis.Redis) (err error) {
	if expire == 0 {
		return nil
	}
	rdb, ctx, _ := redis_conf.Connect()
	data := struct {
		TimeOut  int64
		JsonData string
	}{
		TimeOut:  time.Now().Unix() + int64(expire),
		JsonData: string(json_data),
	}
	json_data2, err := json.Marshal(data)
	if err != nil {
		return err
	}
	expire += uint32(php2go.Rand(60, 600)) // 延长过期时间
	err = rdb.SetEX(ctx, key, json_data2, time.Duration(expire)*time.Second).Err()
	return err
}
