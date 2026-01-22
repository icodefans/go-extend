// 数据缓存
package cache

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// 数据缓存调用
func DataCall(mode Mode, f DataHandlerFunc, data any, args ...any) (err error) {
	// 缓存参数初始化
	init := f(INIT, args...)
	if init.Error != nil {
		return fmt.Errorf("DataCall INIT Error1:%s", init.Error)
	} else if init.Key == "" {
		return fmt.Errorf("DataCall INIT Error2:缓存方法未设置缓存标识")
	} else if init.Redis.Client == nil {
		return fmt.Errorf("DataCall INIT Error3:缓存方法未设置缓存配置")
	} else if init.Data != nil {
		return fmt.Errorf("DataCall INIT Error4:缓存初始化时不能返回数据")
	}
	var argKey string
	if len(args) == 0 {
		// break
	} else if subArgs, err := json.Marshal(args); err != nil {
		return fmt.Errorf("DataCall INIT Error5:缓存方法参数序列化错误，%s", err)
	} else {
		argKey = fmt.Sprintf(":%x", md5.Sum(subArgs))
	}
	key := fmt.Sprintf("data://%s%s", init.Key, argKey)
	expire := time.Duration(1) * time.Millisecond
	if init.Expire > 0 {
		expire = time.Duration(init.Expire) * time.Second
	}
	rdb, ctx, _ := init.Redis.Connect()
	// 缓存数据删除
	if mode == DELETE {
		return rdb.Del(ctx, key).Err()
	}
	// 缓存数据手动设置
	if mode != SET {
		// break
	} else if data == nil {
		return fmt.Errorf("DataCall SET Error1:%s", "缓存数据不能为空")
	} else if jsonData, err := json.Marshal(data); err != nil {
		return fmt.Errorf("DataCall SET Error2:%s", err)
	} else if err := rdb.SetEx(ctx, key, jsonData, expire).Err(); err != nil {
		return fmt.Errorf("DataCall SET Error3:%s", err)
	} else {
		return nil
	}
	// 缓存数据获取，支持缓存关闭
	if value, err := rdb.Get(ctx, key).Result(); err != nil && !errors.Is(err, rdb.KeyNil) {
		return fmt.Errorf("DataCall GET Error1:%s", err)
	} else if init.Expire == 0 || errors.Is(err, rdb.KeyNil) {
		// break
	} else if err := json.Unmarshal([]byte(value), data); err != nil {
		return fmt.Errorf("DataCall GET Error2:%s", err)
	} else {
		return nil
	}
	// 缓存方法调用，缓存不存在则调用方法获取，缓存关闭也需设置缓存，防止后面开启缓存后缓存数据会比较旧
	if result := f(READ, args...); result.Error != nil {
		return fmt.Errorf("DataCall READ Error1:%s", result.Error)
	} else if result.Data == nil {
		return fmt.Errorf("DataCall READ Error2:%s", "缓存数据不能为空")
	} else if jsonData, err := json.Marshal(result.Data); err != nil {
		return fmt.Errorf("DataCall READ Error3:%s", err)
	} else if err := json.Unmarshal(jsonData, data); err != nil {
		return fmt.Errorf("DataCall READ Error4:%s", err)
	} else if err := rdb.SetEx(ctx, key, jsonData, expire).Err(); err != nil {
		return fmt.Errorf("DataCall READ Error5:%s", err)
	}
	return nil
}
