// This file is auto-generated, don't edit it. Thanks.
package aliyun

import (
	"encoding/json"
	"fmt"
	"strings"

	captcha20230305 "github.com/alibabacloud-go/captcha-20230305/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type Aliyun struct {
	AccessKeyId     string `mapstructure:"AccessKeyId"`
	AccessKeySecret string `mapstructure:"AccessKeySecret"`
}

func (aliyun Aliyun) AuthenticateSig(param string) (_err error) {
	// 请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID 和 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例使用环境变量获取 AccessKey 的方式进行调用，仅供参考，建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: &aliyun.AccessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: &aliyun.AccessKeySecret,
	}
	// Endpoint 请参考 https://api.aliyun.com/product/captcha
	config.Endpoint = tea.String("captcha.cn-shanghai.aliyuncs.com")
	client := &captcha20230305.Client{}
	client, _err = captcha20230305.NewClient(config)
	if _err != nil {
		return _err
	}

	verifyIntelligentCaptchaRequest := &captcha20230305.VerifyIntelligentCaptchaRequest{
		CaptchaVerifyParam: tea.String(param),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		_, _err = client.VerifyIntelligentCaptchaWithOptions(verifyIntelligentCaptchaRequest, runtime)
		if _err != nil {
			return _err
		}

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 错误 message
		fmt.Println(tea.StringValue(error.Message))
		// 诊断地址
		var data any
		d := json.NewDecoder(strings.NewReader(tea.StringValue(error.Data)))
		d.Decode(&data)
		// r := reflect.ValueOf(data)
		// fmt.Println("%s", r.MapIndex("Recommend").Interface())
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			return _err
		}
	}
	return _err
}
