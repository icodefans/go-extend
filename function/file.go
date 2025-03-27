package function

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/icodefans/go-extend/define"
	"github.com/syyongx/php2go"
)

// 获取文件MIME信息
func FileMime(filePath string) (string, error) {
	fi, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	// 只需要前 512 个字节就可以了
	buffer := make([]byte, 512)
	_, _ = fi.Read(buffer)
	contentType := http.DetectContentType(buffer)
	return contentType, nil
}

// 复制文件
func FileCopy(src, dst string) (int64, error) {
	// 目录创建
	folderPath := filepath.Dir(dst)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.MkdirAll(folderPath, os.ModePerm) // 0777也可以os.ModePerm
		// os.Chmod(folderPath, 0777)
	}
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// 文件下载
func FileDownLoad(savePath string, fileUrl string) (filePath string, e error) {
	// 文件保存路径设置
	if urlInfo, err := url.ParseRequestURI(fileUrl); err != nil {
		return "", err
	} else if idx := strings.LastIndex(urlInfo.Path, "/"); idx < 0 {
		return "", fmt.Errorf("文件URL地址不正确")
	} else if pathLen := len(urlInfo.Path); pathLen <= idx+1 {
		return "", fmt.Errorf("文件路径不正确")
	} else {
		filePath = fmt.Sprintf("%s/%s", strings.TrimRight(savePath, "/"), php2go.Uniqid(""))
	}
	// 下载文件
	res, err := http.Get(fileUrl)
	if err != nil {
		return "", fmt.Errorf("Http get [%v] failed! %v", fileUrl, err)
	}
	defer res.Body.Close()
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
	if err = os.Rename(filePath, newFilePath); err != nil {
		return "", err
	}
	// 成功返回
	return newFilePath, nil
}

// 检测文件是否存在
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
