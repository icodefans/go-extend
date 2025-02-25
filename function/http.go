package function

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// HTTP请求发起
func HttpCall(method, url string, payload, resData any) (err error) {
	// body参数设置
	var reqBody io.Reader
	if value, ok := payload.([]byte); ok {
		reqBody = strings.NewReader(string(value))
	} else if value, ok := payload.(string); ok {
		reqBody = strings.NewReader(value)
	} else if payload != nil {
		if body, err := json.Marshal(payload); err != nil {
			return err
		} else {
			reqBody = strings.NewReader(string(body))
		}
	}
	// 请求发起
	client := &http.Client{}
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("NewRequestErr %s", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client.DoErr %s", err)
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAllErr %s", err)
	}
	// 结果参数反序列化
	if resData != nil {
		return json.Unmarshal(resBody, &resData)
	}
	return err
}

func HttpFormCall(method, url string, payload, resData any) (err error) {
	var reqBody io.Reader
	if value, ok := payload.([]byte); ok {
		reqBody = strings.NewReader(string(value))
	} else if value, ok := payload.(string); ok {
		reqBody = strings.NewReader(value)
	} else if payload != nil {
		if body, err := json.Marshal(payload); err != nil {
			return err
		} else {
			reqBody = strings.NewReader(string(body))
		}
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return
	}
	req.Header.Add("Cookie", "Accept-Language=zh-cn;Client-Id=103;Client-Version=1.2.1;Client-Channel=huawei;Client-Type=web;Device-UDID=81ed4f86-102d-4e28-bddf-97547048c298;Device-Name=iPhone 13 mini;Device-Resolution=1090*1080;System-Name=Android;System-Version=15.4.1;Team-Name=wahaha;Authorization=1:5776e8530674cdd3d3b2b0d15b1d7623&1724141146:1724054746")
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "api.star-pay.vip")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAllErr %s", err)
	}
	// 结果参数反序列化
	if resData != nil {
		return json.Unmarshal(resBody, &resData)
	}
	return err
}

// HTTP请求发起
func HttpCall2(method, url string, payload, resData any, header map[string]string) (err error) {
	// body参数设置
	var reqBody io.Reader
	if value, ok := payload.([]byte); ok {
		reqBody = strings.NewReader(string(value))
	} else if value, ok := payload.(string); ok {
		reqBody = strings.NewReader(value)
	} else if payload != nil {
		if body, err := json.Marshal(payload); err != nil {
			return err
		} else {
			reqBody = strings.NewReader(string(body))
		}
	}
	// 请求发起
	client := &http.Client{}
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("NewRequestErr %s", err)
	}
	// 请求头设置
	for key, value := range header {
		req.Header.Add(key, value)
	}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client.DoErr %s", err)
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAllErr %s", err)
	}
	// 结果参数反序列化
	if resData != nil {
		return json.Unmarshal(resBody, &resData)
	}
	return err
}
