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
func PagesCall(mode Mode, f PagesHandlerFunc, data any, page *Page, args ...any) (err error) {
	// 验证分页参数
	if page != nil && page.Number == 0 {
		return errors.New("分页编号未设置")
	} else if page != nil && page.Limit == 0 {
		return errors.New("分页大小未设置")
	}
	// 缓存参数初始化
	init := f(INIT, page, args...)
	if init.Error != nil {
		return init.Error
	} else if init.Key == "" {
		return errors.New("缓存方法未设置缓存标识")
	} else if init.Redis.Client == nil {
		return errors.New("缓存方法未设置缓存配置")
	} else if init.Data != nil {
		return errors.New("缓存初始化时不能返回数据")
	}
	var argKey string
	if len(args) > 0 {
		sub_args, err := json.Marshal(args)
		if err != nil {
			return errors.New("缓存方法参数序列化错误")
		}
		argKey = fmt.Sprintf(":%x", md5.Sum(sub_args))
	}
	var path = fmt.Sprintf("page://%s", init.Key)
	key := fmt.Sprintf("%s%s", path, argKey)
	rdb, ctx, _ := init.Redis.Connect()
	// 缓存删除
	if mode == DELETE {
		keys, err := rdb.Keys(ctx, fmt.Sprintf("%s*", path)).Result()
		if err != nil {
			return fmt.Errorf("PagesCache Del Error:%s\n", err)
		} else if len(keys) > 0 {
			rdb.Del(ctx, keys...)
		}
		return nil
	}
	// 缓存获取,支持缓存关闭
	hash_key := "0"
	if page != nil {
		hash_key = fmt.Sprint(page.Number, ":", page.Limit)
	}
	value, err := rdb.HGet(ctx, key, hash_key).Result()
	if err == nil && init.Expire > 0 {
		return json.Unmarshal([]byte(value), data)
	} else if err != nil && !errors.Is(err, rdb.KeyNil) {
		return fmt.Errorf("PagesCache Get Error:%s\n", err)
	}
	// 缓存不存在则调用方法获取
	result := f(READ, page, args...)
	if result.Error != nil {
		return result.Error
	} else if result.Data == nil {
		return errors.New("缓存数据不能为空")
	}
	// 缓存设置
	json_data, err := json.Marshal(result.Data)
	if err != nil {
		return err
	}
	if result.Expire > 0 {
		err = rdb.HSet(ctx, key, hash_key, json_data).Err()
		if err != nil {
			return fmt.Errorf("PagesCache Set Error:%s\n", err)
		}
		err = rdb.Expire(ctx, key, time.Duration(result.Expire)*time.Second).Err()
		if err != nil {
			return fmt.Errorf("PagesCache Exp Error:%s\n", err)
		}
	}
	// 反序列化赋值
	err = json.Unmarshal(json_data, data)
	if err != nil {
		return err
	}
	// 结果返回
	return nil
}
