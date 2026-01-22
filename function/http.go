package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/icodefans/go-extend/define"
	"github.com/syyongx/php2go"
)

// HTTP请求发起
func HttpCall(method, url string, payload, resData any) (err error) {
	// 请求参数设置
	var reqBody io.Reader
	if payload == nil {
		// break
	} else if value, ok := payload.([]byte); ok {
		reqBody = strings.NewReader(string(value))
	} else if value, ok := payload.(string); ok {
		reqBody = strings.NewReader(value)
	} else if body, err := json.Marshal(payload); err != nil {
		return fmt.Errorf("HttpCall json.MarshalErr %s", err)
	} else {
		reqBody = strings.NewReader(string(body))
	}
	// 请求发起
	var res *http.Response
	if req, err := http.NewRequest(method, url, reqBody); err != nil {
		return fmt.Errorf("HttpCall NewRequestErr %s", err)
	} else if res, err = (&http.Client{}).Do(req); err != nil {
		return fmt.Errorf("HttpCall client.DoErr %s", err)
	}
	defer func(Body io.ReadCloser) { err = Body.Close() }(res.Body)
	// 请求结果获取
	if resBody, err := io.ReadAll(res.Body); err != nil {
		return fmt.Errorf("HttpCall io.ReadAllErr %s", err)
	} else if resData == nil {
		// break
	} else if err := json.Unmarshal(resBody, &resData); err != nil {
		return fmt.Errorf("HttpCall json.UnmarshalErr %s", err)
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

// 文件上传
func HttpFileUpload(uploadURL, filePath string, resData any) error {
	// 打开本地文件
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 获取文件信息用于表单
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	filename := fileInfo.Name()

	// 创建缓冲写入器用于构造multipart表单
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 创建文件字段并写入文件内容
	fileField, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return err
	}
	_, err = io.Copy(fileField, file)
	if err != nil {
		return err
	}

	// 必须调用Close()来写入表单结束符
	writer.Close()

	// 创建HTTP请求
	req, err := http.NewRequest("POST", uploadURL, &buf)
	if err != nil {
		return err
	}
	// 设置Content-Type为multipart/form-data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) { _ = Body.Close() }(res.Body)

	// 检查响应状态
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("服务器错误: %s", res.Status)
	}
	// 请求结果获取
	if resBody, err := io.ReadAll(res.Body); err != nil {
		return fmt.Errorf("HttpCall io.ReadAllErr %s", err)
	} else if resData == nil {
		// break
	} else if err := json.Unmarshal(resBody, &resData); err != nil {
		return fmt.Errorf("HttpCall json.UnmarshalErr %s", err)
	}
	return nil
}

// 文件下载
func HttpFileDownLoad(fileUrl, savePath string, resData any) (filePath string, err error) {
	// 文件地址验证
	if urlInfo, err := url.ParseRequestURI(fileUrl); err != nil {
		return "", err
	} else if idx := strings.LastIndex(urlInfo.Path, "/"); idx < 0 {
		return "", fmt.Errorf("文件URL地址不正确")
	} else if pathLen := len(urlInfo.Path); pathLen <= idx+1 {
		return "", fmt.Errorf("文件路径不正确")
	}
	// 文件保存路径设置
	{
		filePath = fmt.Sprintf("%s/%s", strings.TrimRight(savePath, "/"), php2go.Uniqid(""))
	}
	// 请求发起
	var res *http.Response
	if req, err := http.NewRequest("GET", fileUrl, nil); err != nil {
		return "", err
	} else if res, err = (&http.Client{}).Do(req); err != nil {
		return "", err
	}
	defer func(body io.ReadCloser) { _ = body.Close() }(res.Body)
	if res.StatusCode != 200 {
		return "", fmt.Errorf("文件下载失败")
	}
	// 响应头非文件类型，则返回结果内容
	if matches := regexp.MustCompile(
		fmt.Sprintf(`(%s)`, strings.Join([]string{"application/json", "text/html"}, "|")),
	).FindAllString(res.Header.Get("Content-Type"), -1); len(matches) == 0 {
		// break
	} else if resBody, err := io.ReadAll(res.Body); err != nil {
		return "", fmt.Errorf("HttpCall io.ReadAllErr %s", err)
	} else if err := json.Unmarshal(resBody, &resData); err != nil {
		return "", fmt.Errorf("HttpCall json.UnmarshalErr %s", err)
	} else {
		return "", nil
	}
	// 文件保存
	if err := os.MkdirAll(savePath, 0755); err != nil {
		return "", fmt.Errorf("目录创建失败: %v", err)
	} else if file, err := os.Create(filePath); err != nil {
		return "", fmt.Errorf("文件创建出错: %v", err)
	} else if _, err = io.Copy(file, res.Body); err != nil {
		return "", fmt.Errorf("文件保存出错: %v", err)
	}
	// 通过文件mime信息获取文件后缀
	var ext string
	if mime, err := FileMime(filePath); err != nil {
		return "", err
	} else if ext, _ = define.MIME_EXT[mime]; ext == "" {
		return "", fmt.Errorf("MIME(%s)信息文件后缀未定义", mime)
	}
	newFilePath := fmt.Sprintf("%s%s", filePath, ext)
	// 修改文件名
	if err := os.Rename(filePath, newFilePath); err != nil {
		return "", err
	}
	// 成功返回
	return newFilePath, nil
}
