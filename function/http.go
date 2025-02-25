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
