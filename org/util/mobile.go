package util

import (
	"fmt"
	"strconv"
	"strings"
)

// 国际手机号
type Mobile struct {
	Region uint32 `json:"region" lable:"国际区号"`
	Number uint64 `json:"number" lable:"手机号"`
}

// 国际手机号解析
// 手机号码国际书写方式 +国家代码-对方用户手机号码 +86-18876543210
func MobileParse(mobile string) (obj *Mobile, err error) {
	// 参数验证
	var mobileArr []string
	if mobile == "" {
		return nil, fmt.Errorf("手机号不能为空")
	} else if mobile[0:1] != "+" {
		return nil, fmt.Errorf("手机号需要加号开头")
	} else if mobileArr = strings.Split(mobile, "-"); len(mobileArr) != 2 {
		return nil, fmt.Errorf("手机区号与号码需要剪号分隔")
	}
	// 类型转换
	if region, err := strconv.Atoi(mobileArr[0][1:]); err != nil {
		return nil, err
	} else if number, err := strconv.Atoi(mobileArr[1]); err != nil {
		return nil, err
	} else {
		obj = &Mobile{Region: uint32(region), Number: uint64(number)}
	}
	// 成功返回
	return obj, nil
}
