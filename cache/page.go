// 分页缓存
package cache

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// 分页缓存调用
func PageCall(mode Mode, f PageHandlerFunc, subKey string, data any, page *Page, args ...any) (err error) {
	// 验证分页参数
	if page != nil && page.Number <= 0 {
		return fmt.Errorf("PageCall PAGE Error1:分页编号未设置")
	} else if page != nil && page.Limit <= 0 {
		return fmt.Errorf("PageCall PAGE Error2:分页大小未设置")
	}
	// 缓存参数初始化
	init := f(INIT, subKey, page, args...)
	if init.Error != nil {
		return fmt.Errorf("PageCall INIT Error1:%s", init.Error)
	} else if init.Key == "" {
		return fmt.Errorf("PageCall INIT Error2:缓存方法未设置缓存标识")
	} else if init.Redis.Client == nil {
		return fmt.Errorf("PageCall INIT Error3:缓存方法未设置缓存配置")
	} else if init.Data != nil {
		return fmt.Errorf("PageCall INIT Error4:缓存初始化时不能返回数据")
	}
	var argKey string
	if len(args) == 0 {
		// break
	} else if argsJson, err := json.Marshal(args); err != nil {
		return fmt.Errorf("PageCall INIT Error5:缓存方法参数序列化错误，%s", err)
	} else {
		argKey = fmt.Sprintf(":%x", md5.Sum(argsJson))
	}
	path := fmt.Sprintf("page://%s", init.Key)
	key := fmt.Sprintf("%s%s", path, argKey)
	expire := time.Duration(1) * time.Millisecond
	if init.Expire > 0 {
		expire = time.Duration(init.Expire) * time.Second
	}
	hashKey := "0:0"
	if page != nil {
		hashKey = fmt.Sprintf("%d:%d", page.Number, page.Limit)
	}
	rdb, ctx, _ := init.Redis.Connect()
	// 缓存数据删除
	if mode != DELETE {
		// break
	} else if _, err := rdb.DelX(ctx, fmt.Sprintf("%s*", path)); err != nil {
		return fmt.Errorf("PageCall DELETE Error1:%s", err)
	} else {
		return nil
	}
	// 缓存数据获取，支持缓存关闭
	if jsonData, err := rdb.HGet(ctx, key, hashKey).Result(); err != nil && !errors.Is(err, rdb.KeyNil) {
		return fmt.Errorf("PageCall GET Error1:%s", err)
	} else if init.Expire == 0 || errors.Is(err, rdb.KeyNil) {
		// break
	} else if err := json.Unmarshal([]byte(jsonData), data); err != nil {
		return fmt.Errorf("PageCall GET Error2:%s", err)
	} else {
		return nil
	}
	// 缓存方法调用，调用成功设置缓存数据
	if result := f(READ, subKey, page, args...); result.Error != nil {
		return fmt.Errorf("PageCall SET Error1:%s", result.Error)
	} else if result.Data == nil {
		return fmt.Errorf("PageCall SET Error2:%s", "缓存数据不能为空")
	} else if jsonData, err := json.Marshal(result.Data); err != nil {
		return fmt.Errorf("PageCall SET Error3:%s", err)
	} else if err := rdb.HSet(ctx, key, hashKey, jsonData).Err(); err != nil {
		return fmt.Errorf("PageCall SET Error4:%s", err)
	} else if err := rdb.Expire(ctx, key, expire).Err(); err != nil {
		return fmt.Errorf("PageCall SET Error5:%s", err)
	} else if err := json.Unmarshal(jsonData, data); err != nil {
		return fmt.Errorf("PageCall SET Error6:%s", err)
	}
	return nil
}
