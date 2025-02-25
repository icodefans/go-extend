package storage

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/syyongx/php2go"
)

// 阿里云OOS配置
type Aliyun struct {
	AccessKeyId     string `mapstructure:"AccessKeyId"`
	AccessKeySecret string `mapstructure:"AccessKeySecret"`
	Endpoint        string `mapstructure:"Endpoint"`
	BucketName      string `mapstructure:"BucketName"`
	Root            string `mapstructure:"Root"`
	Domain          string `mapstructure:"Domain"`
	Bucket          *oss.Bucket
}

// 实例化客户端
func (aliyun *Aliyun) Client() *Aliyun {
	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	client, err := oss.New(aliyun.Endpoint, aliyun.AccessKeyId, aliyun.AccessKeySecret)
	if err != nil {
		panic(err)
	}

	// 填写存储空间名称，例如examplebucket。
	bucket, err := client.Bucket(aliyun.BucketName)
	if err != nil {
		panic(err)
	}
	aliyun.Bucket = bucket
	return aliyun
}

// 上传本地文件
// @param root string 保存至目录
// @param filePath 本地文件路径
// @return string,error
func (aliyun Aliyun) UploadLoclFile(savePath, filePath string) (objectKey string, domain string, err error) {
	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	client, err := oss.New(aliyun.Endpoint, aliyun.AccessKeyId, aliyun.AccessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		return "", "", err
	}

	// 填写存储空间名称，例如examplebucket。
	bucket, err := client.Bucket(aliyun.BucketName)
	if err != nil {
		fmt.Println("Error:", err)
		return "", "", err
	}

	var saveDir string
	if aliyun.Root != "" {
		saveDir += strings.Trim(aliyun.Root, "/") + "/"
	}
	if savePath != "" {
		saveDir += strings.Trim(savePath, "/") + "/"
	}
	fileName := php2go.Uniqid("") + strings.ToLower(filepath.Ext(filePath))
	objectKey = fmt.Sprint(saveDir, fileName)

	// 开始文件上传
	err = bucket.PutObjectFromFile(objectKey, filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return "", "", err
	}
	// 成功返回
	return objectKey, strings.Trim(aliyun.Domain, "/") + "/", nil
}

// 下载文件至指定位置
// @param objectKey string 对象文件
// @param tempFile 临时文件
// @return error
func (aliyun Aliyun) DownloadLocalFile(objectKey, tempFile string) (err error) {
	err = aliyun.Bucket.GetObjectToFile(objectKey, tempFile)
	// 成功返回
	return
}
