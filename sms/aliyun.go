package sms

import (
	"errors"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	aliyunUtil "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
)

type Aliyun struct {
	AccessKeyId     string `mapstructure:"AccessKeyId"`
	AccessKeySecret string `mapstructure:"AccessKeySecret"`
	SignName        string `mapstructure:"SignName"`     // 短信签名
	TemplateCode    string `mapstructure:"TemplateCode"` // 模板标识
}

// 客户端创建
func (*Aliyun) client(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

// 短信发送
func (aliyun *Aliyun) send(req dysmsapi20170525.SendSmsRequest) (_err error) {
	client, _err := aliyun.client(tea.String(aliyun.AccessKeyId), tea.String(aliyun.AccessKeySecret))
	if _err != nil {
		return _err
	}

	defer func() {
		if r := tea.Recover(recover()); r != nil {
			_err = r
		}
	}()

	runtime := &aliyunUtil.RuntimeOptions{}
	result, _err := client.SendSmsWithOptions(&req, runtime)
	if _err != nil {
		return _err
	}

	if *result.Body.Code != "OK" {
		_err = errors.New(result.String())
		return
	}

	return _err
}

// 验证码短信发送
func (aliyun *Aliyun) CaptchaSmsSend(phoneNumber, code, template string) (err error) {
	req := dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String(aliyun.SignName),
		TemplateCode:  tea.String(template),
		PhoneNumbers:  tea.String(phoneNumber),
		TemplateParam: tea.String(`{"code":"` + code + `"}`),
	}
	return aliyun.send(req)
}
