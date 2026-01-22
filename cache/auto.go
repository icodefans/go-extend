// 自动缓存
package cache

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/icodefans/go-extend/database/redis"
	"github.com/syyongx/php2go"
)

// 自动缓存调用
func AutoCall(mode Mode, f DataHandlerFunc, data any, args ...any) (err error) {
	// 缓存参数初始化
	init := f(INIT, args...)
	if init.Error != nil {
		return fmt.Errorf("AutoCall INIT Error1:%s", init.Error)
	} else if init.Key == "" {
		return fmt.Errorf("AutoCall INIT Error2:缓存方法未设置缓存标识")
	} else if init.Redis.Client == nil {
		return fmt.Errorf("AutoCall INIT Error3:缓存方法未设置缓存配置")
	} else if init.Data != nil {
		return fmt.Errorf("AutoCall INIT Error4:缓存初始化时不能返回数据")
	}
	var argKey string
	if len(args) == 0 {
		// break
	} else if subArgs, err := json.Marshal(args); err != nil {
		return fmt.Errorf("AutoCall INIT Error5:缓存方法参数序列化错误%s", err)
	} else {
		argKey = fmt.Sprintf(":%x", md5.Sum(subArgs))
	}
	key := fmt.Sprintf("auto://%s%s", init.Key, argKey)
	lock := fmt.Sprint(key, ".lock")
	rdb, ctx, _ := init.Redis.Connect()
	// 缓存数据删除
	if mode == DELETE {
		return rdb.Del(ctx, key).Err()
	}
	// 缓存数据手动设置
	if mode != SET {
		// break
	} else if data == nil {
		return fmt.Errorf("AutoCall SET Error1:%s", "缓存数据不能为空")
	} else if _, err := autoSet(key, lock, &data, init.Expire, init.Redis); err != nil {
		return fmt.Errorf("AutoCall SET Error2:%s", err)
	} else {
		return nil
	}
	// 缓存数据获取，支持缓存关闭
	if jsonData, err := autoGet(key, lock, 0, init.Redis); err != nil && !errors.Is(err, rdb.KeyNil) {
		return fmt.Errorf("AutoCall GET Error1:%s", err)
	} else if init.Expire == 0 || errors.Is(err, rdb.KeyNil) {
		// break
	} else if err := json.Unmarshal([]byte(jsonData), data); err != nil {
		return fmt.Errorf("AutoCall GET Error2:key.lock:%s,%s", rdb.Del(ctx, lock).Name(), err)
	} else {
		return nil
	}
	// 缓存方法调用，调用成功则设置缓存数据
	if result := f(READ, args...); result.Error != nil {
		return fmt.Errorf("AutoCall READ Error1:key.lock:%s,%s", rdb.Del(ctx, lock).Name(), result.Error)
	} else if result.Data == nil {
		return fmt.Errorf("AutoCall READ Error2:key.lock:%s,%s", rdb.Del(ctx, lock).Name(), "缓存数据不能为空")
	} else if jsonData, err := autoSet(key, lock, &result.Data, result.Expire, result.Redis); err != nil {
		return fmt.Errorf("AutoCall READ Error3:key.lock:%s,%s", rdb.Del(ctx, lock).Name(), err)
	} else if err = json.Unmarshal(jsonData, data); err != nil {
		return fmt.Errorf("AutoCall READ Error4:%s", err)
	}
	return nil
}

// 自动缓存获取
func autoGet(key, lock string, count uint32, redis redis.Redis) (jsonData string, err error) {
	var (
		lockExpire = uint32(30) // 单例写入锁缓存时间
	)
	// 缓存读取
	rdb, ctx, _ := redis.Connect()
	if jsonData, err = rdb.Get(ctx, key).Result(); err == nil {
		// break
	} else if err != nil && !errors.Is(err, rdb.KeyNil) {
		return "", fmt.Errorf("autoGet.rdb.Get:%s", err) // 缓存获取报错
	} else if incrV, err := rdb.IncrX(ctx, lock, lockExpire); err != nil {
		return "", fmt.Errorf("autoGet.rdb.IncrX1:%s", err) // 自增锁设置报错
	} else if incrV == 1 {
		return "", rdb.KeyNil // 缓存不存在处理逻辑（让一个请求去生成缓存，其他请求等待或报错）
	} else if time.Sleep(time.Second * 1); incrV > 1 && count < 60 {
		return autoGet(key, lock, count+1, redis) // 等待缓存生成，每间隔1秒递归调用获取缓存，上限60次
	} else if log.Printf("AutoGetErr:%s,incrV:%d", key, incrV); true {
		return "", fmt.Errorf("autoGet.缓存数据载入中，请稍后再试~") // 超过60秒没有获取到缓存，缓存获取失败
	}
	// 缓存读取成功，解析缓存
	var data struct {
		TimeOut  int64
		JsonData string
	}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return "", fmt.Errorf("autoGet.json.Unmarshal:%s", err)
	}
	// 缓存未过期，直接返回
	if time.Now().Unix() <= data.TimeOut {
		return data.JsonData, nil
	}
	// 缓存过期处理逻辑（让一个请求去生成缓存，其他请求读取旧的缓存）
	if incrV, err := rdb.IncrX(ctx, lock, lockExpire); err != nil {
		return "", fmt.Errorf("autoGet.rdb.IncrX2:%s", err)
	} else if incrV == 1 {
		return "", rdb.KeyNil
	}
	// 成功返回（读取旧的缓存）
	return data.JsonData, nil
}

// 自动缓存设置
// 缓存关闭也需设置缓存，防止后面开启缓存后缓存数据会比较旧
func autoSet(key, lock string, resultData any, expire uint32, redis redis.Redis) (jsonData []byte, err error) {
	if jsonData, err = json.Marshal(resultData); err != nil {
		return nil, fmt.Errorf("autoSet.json.Marshal1:%s", err)
	}
	var data struct {
		TimeOut  int64
		JsonData string
	}
	{
		data.TimeOut = time.Now().Unix() + int64(expire)
		data.JsonData = string(jsonData)
	}
	expire += uint32(php2go.Rand(60, 600)) // 过期时间随机延长，防止缓存集中到期发生雪崩
	rdb, ctx, _ := redis.Connect()
	if jsonDataNew, err := json.Marshal(data); err != nil {
		return nil, fmt.Errorf("autoSet.json.Marshal2:%s", err)
	} else if err := rdb.SetEx(ctx, key, jsonDataNew, time.Duration(expire)*time.Second).Err(); err != nil {
		return nil, fmt.Errorf("autoSet.rdb.SetEX:%s", err)
	}
	// 缓存设置完成，删除单例写入锁
	return jsonData, rdb.Del(ctx, lock).Err()
}
