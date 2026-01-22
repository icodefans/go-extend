package method

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/icodefans/go-extend/define"
	"github.com/icodefans/go-extend/function"
	"github.com/icodefans/go-extend/validate"
	"github.com/syyongx/php2go"
)

// HTTP请求模块
type Http struct {
	Url     string            `json:"url" validate:"required,url" label:"地址"`
	Method  string            `json:"method" validate:"required,oneof=GET POST PUT DELETE" label:"请求方式"`
	Header  map[string]string `json:"header" validate:"omitempty" label:"请求头"`
	Timeout uint32            `json:"timeout" validate:"omitempty" label:"超时时间,单位秒"`
}

// 请求发起
func (h *Http) Call(payload, resData any) (err error) {
	// 参数验证
	if err = validate.Validate(h); err != nil {
		return err
	}
	// 请求body参数设置
	var reqBody io.Reader
	if payload == nil {
		// break
	} else if body, ok := payload.([]byte); ok {
		reqBody = strings.NewReader(string(body))
	} else if body, ok := payload.(string); ok {
		reqBody = strings.NewReader(body)
	} else if body, err := json.Marshal(payload); err != nil {
		return fmt.Errorf("HttpCall json.MarshalErr %s", err)
	} else {
		reqBody = strings.NewReader(string(body))
	}
	// 客户端设置
	var res *http.Response
	var req *http.Request
	var client = &http.Client{
		Timeout: time.Duration(h.Timeout) * time.Second,
	}
	if req, err = http.NewRequest(h.Method, h.Url, reqBody); err != nil {
		return fmt.Errorf("HttpCall NewRequestErr %s", err)
	}
	// 请求头设置
	for key, value := range h.Header {
		req.Header.Add(key, value)
	}
	// 请求发起
	if res, err = client.Do(req); err != nil {
		return fmt.Errorf("HttpCall client.DoErr %s", err)
	} else if res.StatusCode >= 400 {
		return fmt.Errorf("HttpCall client.DoErr %s", res.Status)
	}
	// 请求结果获取
	defer func(Body io.ReadCloser) { err = Body.Close() }(res.Body)
	if resData == nil {
		// break
	} else if resBody, err := io.ReadAll(res.Body); err != nil {
		return fmt.Errorf("HttpCall io.ReadAllErr %s", err)
	} else if err := json.Unmarshal(resBody, &resData); err != nil {
		return fmt.Errorf("HttpCall json.UnmarshalErr %s", err)
	}
	return err
}

// 文件上传
func (h *Http) Upload(filePath string, payload, resData any) (err error) {
	// 参数验证
	if err = validate.Validate(h); err != nil {
		return err
	}
	// 打开本地文件
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) { _ = file.Close() }(file)
	// 创建缓冲写入器用于构造multipart表单
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	// 获取文件信息用于表单,  创建文件字段并写入文件内容
	if fileInfo, err := file.Stat(); err != nil {
		return err
	} else if fileField, err := writer.CreateFormFile("file", fileInfo.Name()); err != nil {
		return err
	} else if _, err = io.Copy(fileField, file); err != nil {
		return err
	}
	// 必须调用Close()来写入表单结束符
	if err = writer.Close(); err != nil {
		return err
	}
	// // 请求body参数设置
	// var reqBody io.Reader
	// if payload == nil {
	// 	// break
	// } else if body, ok := payload.([]byte); ok {
	// 	reqBody = strings.NewReader(string(body))
	// } else if body, ok := payload.(string); ok {
	// 	reqBody = strings.NewReader(body)
	// } else if body, err := json.Marshal(payload); err != nil {
	// 	return fmt.Errorf("HttpCall json.MarshalErr %s", err)
	// } else {
	// 	reqBody = strings.NewReader(string(body))
	// }
	// 客户端设置
	var res *http.Response
	var req *http.Request
	var client = &http.Client{
		Timeout: time.Duration(h.Timeout) * time.Second,
	}
	// 创建HTTP请求
	if req, err = http.NewRequest(h.Method, h.Url, &buf); err != nil {
		return err
	}
	// 设置Content-Type为multipart/form-data
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// 请求头设置
	for key, value := range h.Header {
		req.Header.Add(key, value)
	}
	// 发送请求
	if res, err = client.Do(req); err != nil {
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
func (h *Http) DownLoad(savePath string, payload, resData any) (filePath *string, err error) {
	// 参数验证
	if err = validate.Validate(h); err != nil {
		return nil, err
	}
	// 请求body参数设置
	var reqBody io.Reader
	if payload == nil {
		// break
	} else if body, ok := payload.([]byte); ok {
		reqBody = strings.NewReader(string(body))
	} else if body, ok := payload.(string); ok {
		reqBody = strings.NewReader(body)
	} else if body, err := json.Marshal(payload); err != nil {
		return nil, fmt.Errorf("HttpCall json.MarshalErr %s", err)
	} else {
		reqBody = strings.NewReader(string(body))
	}
	// 客户端设置
	var res *http.Response
	var req *http.Request
	var client = &http.Client{
		Timeout: time.Duration(h.Timeout) * time.Second,
	}
	if req, err = http.NewRequest(h.Method, h.Url, reqBody); err != nil {
		return nil, fmt.Errorf("HttpCall NewRequestErr %s", err)
	}
	// 请求头设置
	for key, value := range h.Header {
		req.Header.Add(key, value)
	}
	// 请求发起
	if res, err = client.Do(req); err != nil {
		return nil, fmt.Errorf("HttpCall client.DoErr %s", err)
	}
	defer func(Body io.ReadCloser) { err = Body.Close() }(res.Body)
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("文件下载失败")
	}
	// 响应头非文件类型，则返回结果内容
	if matches := regexp.MustCompile(
		fmt.Sprintf(`(%s)`, strings.Join([]string{"application/json", "text/html"}, "|")),
	).FindAllString(res.Header.Get("Content-Type"), -1); len(matches) == 0 {
		// break
	} else if resBody, err := io.ReadAll(res.Body); err != nil {
		return nil, fmt.Errorf("HttpCall io.ReadAllErr %s", err)
	} else if err := json.Unmarshal(resBody, &resData); err != nil {
		return nil, fmt.Errorf("HttpCall json.UnmarshalErr %s", err)
	} else {
		return nil, nil
	}
	// 文件保存
	filePath = new(string)
	*filePath = fmt.Sprintf("%s/%s", strings.TrimRight(savePath, "/"), php2go.Uniqid(""))
	if err := os.MkdirAll(savePath, 0755); err != nil {
		return nil, fmt.Errorf("目录创建失败: %v", err)
	} else if file, err := os.Create(*filePath); err != nil {
		return nil, fmt.Errorf("文件创建出错: %v", err)
	} else if _, err = io.Copy(file, res.Body); err != nil {
		return nil, fmt.Errorf("文件保存出错: %v", err)
	}
	// 通过文件mime信息获取文件后缀
	var ext string
	if ext = function.FileGetExt(h.Url); ext != "" {
		// next
	} else if mime, err := function.HttpFileMime(*filePath); err != nil {
		return nil, err
	} else if ext, _ = define.MIME_EXT[mime]; ext == "" {
		return nil, fmt.Errorf("MIME(%s)信息文件后缀未定义", mime)
	}
	newFilePath := fmt.Sprintf("%s%s", *filePath, ext)
	// 修改文件名
	if err := os.Rename(*filePath, newFilePath); err != nil {
		return nil, err
	}
	// 成功返回
	return &newFilePath, nil
}
