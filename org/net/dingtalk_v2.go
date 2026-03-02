package net

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DingTalk_v2 struct {
	Webhook string
}

// SendDingTalkMessage 发送钉钉消息并@指定人
func (dingtalk DingTalk_v2) Markdown(title, text string, atMobiles []string, isAtAll bool) error {
	// DingTalkMessage 钉钉机器人消息结构体
	var msg struct {
		MsgType  string `json:"msgtype"`
		Markdown struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		} `json:"markdown"`
		At struct {
			AtMobiles []string `json:"atMobiles"` // 要@的手机号列表
			AtUserIds []string `json:"atUserIds"` // 要@的userId列表（二选一即可）
			IsAtAll   bool     `json:"isAtAll"`   // 是否@所有人
		} `json:"at"`
	}
	{
		msg.MsgType = "markdown"
		msg.Markdown.Title = title
		msg.Markdown.Text = text
		msg.At.AtMobiles = atMobiles
		msg.At.IsAtAll = isAtAll
	}

	// 将结构体转换为JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}

	// 发送POST请求
	resp, err := http.Post(dingtalk.Webhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	return nil
}

// SendDingTalkMessage 发送钉钉消息并@指定人
func (dingtalk DingTalk_v2) Text(content string, atMobiles []string, isAtAll bool) error {
	// DingTalkMessage 钉钉机器人消息结构体
	var msg struct {
		MsgType string `json:"msgtype"` // 消息类型，固定为text
		Text    struct {
			Content string `json:"content"` // 消息内容
		} `json:"text"`
		At struct {
			AtMobiles []string `json:"atMobiles"` // 要@的手机号列表
			AtUserIds []string `json:"atUserIds"` // 要@的userId列表（二选一即可）
			IsAtAll   bool     `json:"isAtAll"`   // 是否@所有人
		} `json:"at"`
	}
	{
		msg.MsgType = "text"
		msg.Text.Content = content
		msg.At.AtMobiles = atMobiles
		msg.At.IsAtAll = isAtAll
	}

	// 将结构体转换为JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}

	// 发送POST请求
	resp, err := http.Post(dingtalk.Webhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	return nil
}
