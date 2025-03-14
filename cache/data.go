// 数据缓存
package cache

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"time"
)

// 数据缓存调用
func DataCall(mode Mode, f DataHandlerFunc, data any, args ...any) (err error) {
	// 异常处理
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		rec := recover()
		if rec != nil {
			err = fmt.Errorf("数据缓存调用异常")
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
	key := fmt.Sprint("data://", init.Key, sub_key)
	rdb, ctx, _ := init.Redis.Connect()
	// 缓存删除
	if mode == DELETE {
		rdb.Del(ctx, key)
		return nil
	}
	// 缓存手动设置
	if mode == SET && data == nil {
		return errors.New("缓存数据不能为空1")
	} else if mode == SET {
		if init.Expire == 0 {
			return nil
		} else if json_data, err := json.Marshal(data); err != nil {
			return fmt.Errorf("DataCache JsonEncode1 Error:%s\n", err)
		} else if err = rdb.SetEX(ctx, key, json_data, time.Duration(init.Expire)*time.Second).Err(); err != nil {
			return fmt.Errorf("DataCache Set Error:%s\n", err)
		}
		return nil
	}
	// 缓存获取,支持缓存关闭
	if value, err := rdb.Get(ctx, key).Result(); err == nil && init.Expire > 0 {
		if err := json.Unmarshal([]byte(value), data); err != nil {
			return fmt.Errorf("DataCache JsonDecode1 Error:%s\n", err)
		}
		return nil
	} else if err != nil && !errors.Is(err, rdb.KeyNil) {
		return fmt.Errorf("DataCache Get Error:%s\n", err)
	}
	// 缓存自动设置，缓存不存在则调用方法获取
	if result := f(READ, args...); result.Error != nil {
		return result.Error
	} else if result.Data == nil {
		return errors.New("缓存数据不能为空2")
	} else if json_data, err := json.Marshal(result.Data); err != nil {
		return fmt.Errorf("DataCache JsonEncode1 Error:%s\n", err)
	} else if err = json.Unmarshal(json_data, data); err != nil {
		return fmt.Errorf("DataCache JsonDecode2 Error:%s\n", err)
	} else if result.Expire > 0 {
		if err = rdb.SetEX(ctx, key, json_data, time.Duration(result.Expire)*time.Second).Err(); err != nil {
			return fmt.Errorf("DataCache Set Error:%s\n", err)
		}
	}
	return nil
}
