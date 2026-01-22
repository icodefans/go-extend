package storage

import (
	"fmt"
	"strings"
	"time"

	"github.com/icodefans/go-extend/function"
	"github.com/syyongx/php2go"
)

// 本地存储
type Disk struct {
	BucketPath string `mapstructure:"BucketPath"`
	Root       string `mapstructure:"Root"`
	Domain     string `mapstructure:"Domain"`
}

// 实例化客户端
func (aliyun *Disk) Client() *Disk {
	return aliyun
}

// 上传本地文件
// @param root string 保存至目录
// @param filePath 本地文件路径
// @return string,error
func (aliyun *Disk) UploadLoclFile(savePath, filePath string) (objectKey string, domain string, err error) {
	var (
		BucketPath string
		saveRoot   string
	)
	if aliyun.BucketPath != "" {
		BucketPath = strings.TrimRight(aliyun.BucketPath, "/")
	}
	if aliyun.Root != "" {
		saveRoot = strings.Trim(aliyun.Root, "/") + "/"
	}
	if savePath != "" {
		savePath = strings.Trim(savePath, "/")
	}

	// 保存目录增加日期
	DateDir := time.Now().Format("20060102")
	savePath = fmt.Sprintf("%s%s/%s/", saveRoot, savePath, DateDir)
	BucketPath = fmt.Sprintf("%s/%s", BucketPath, savePath)

	fileName := php2go.Uniqid("") + strings.ToLower(function.FileGetExt(filePath))
	objectKey = fmt.Sprint(savePath, fileName)
	copyPath := fmt.Sprint(BucketPath, fileName)

	// 复制文件
	if _, err := function.FileCopy(filePath, copyPath); err != nil {
		return "", "", err
	}
	// 成功返回
	return objectKey, strings.Trim(aliyun.Domain, "/") + "/", nil
}

// 下载文件至指定位置
// @param objectKey string 对象文件
// @param tempFile 临时文件
// @return error
func (aliyun *Disk) DownloadLocalFile(objectKey, tempFile string) error {
	var (
		BucketPath string
	)
	if aliyun.Root != "" {
		BucketPath = strings.TrimRight(aliyun.Root, "/")
	}

	// 源文件路径
	filePath := fmt.Sprintf("%s/%s", BucketPath, objectKey)

	// 复制文件
	if _, err := function.FileCopy(filePath, tempFile); err != nil {
		return err
	}
	return nil
}
