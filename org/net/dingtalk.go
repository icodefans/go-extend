package net

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	_ "github.com/blinkbean/dingtalk"
)

type DingTalk struct {
	Secret string `mapstructure:"Secret"`
	Token  string `mapstructure:"Token"`
}

type markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
type at struct {
	AtMobiles string `json:"atMobiles"`
	IsAtAll   bool   `json:"isAtAll"`
}
type dingTalkData struct {
	MsgType  string   `json:"msgtype"`
	Markdown markdown `json:"markdown"`
	At       at       `json:"at"`
}

/** markdown
 * @param  string  $hook
 * @param  string  $title
 * @param  string  $text
 * @param  array   $mobile
 * @param  bool    $at_all
 * @return bool|string
 */
func (config *DingTalk) Markdown(title, text, atMobiles string, isAtAll bool) {
	// 组织数据
	data := dingTalkData{
		MsgType: "markdown",
		Markdown: markdown{
			Title: title,
			Text:  text,
		},
		At: at{
			AtMobiles: atMobiles,
			IsAtAll:   isAtAll,
		},
	}
	json_data, _ := json.Marshal(data)

	// 签名字符串
	timestamp := time.Now().UnixNano() / 1e6
	sign := sign(timestamp, config.Secret)

	// 请求发送
	url := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s&timestamp=%d&sign=%s", config.Token, timestamp, sign)
	curl(url, string(json_data))
}

// 请求参数签名
func sign(t int64, secret string) string {
	strToHash := fmt.Sprintf("%d\n%s", t, secret)
	hmac256 := hmac.New(sha256.New, []byte(secret))
	hmac256.Write([]byte(strToHash))
	data := hmac256.Sum(nil)
	return base64.StdEncoding.EncodeToString(data)
}

/** curl
 * @param string url
 * @param string json
 * @return string,error
 */
func curl(url, json string) (string, error) {
	method := "POST"

	payload := strings.NewReader(json)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
