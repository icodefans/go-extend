package net

import (
	"fmt"

	"github.com/pquerna/otp/totp"
)

// 谷歌身份认证
type GoogleAuthenticator struct {
	Issuer      string `json:"issuer" label:"机构名称"`
	AccountName string `json:"account_name" label:"账号名称"`
}

// 生成一个随机的密钥
func (g *GoogleAuthenticator) Generate() (secretKey, urlKey string, err error) {
	// 生成一个随机的密钥
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      g.Issuer,      // 机构名称
		AccountName: g.AccountName, // 账号名称
	})
	if err != nil {
		return
	}
	// 生成的信息可以存储到mysql 或者redis 里面方便下次验证
	secretKey = key.Secret() // 生成的 密钥
	urlKey = key.URL()       // 生成的二维码链接 ,可以通过链接进行二维码绑定到应用上
	return
}

// 验证OTP码code 是用户输入的
func (g *GoogleAuthenticator) Validate(code, secretKey string) (ok bool) {
	return totp.Validate(code, secretKey)
}

// 生成的二维码链接 ,可以通过链接进行二维码绑定到应用上
func (g *GoogleAuthenticator) UrlKey(secretKey string) string {
	return fmt.Sprintf("otpauth://totp/%s:%s?algorithm=SHA1&digits=6&issuer=%s&period=30&secret=%s",
		g.Issuer,
		g.AccountName,
		g.Issuer,
		secretKey,
	)
}
