package storage

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/syyongx/php2go"
)

// AWS S3配置
type Aws struct {
	AccessKeyId     string `mapstructure:"AccessKeyId"`
	SecretAccesskey string `mapstructure:"SecretAccesskey"`
	Region          string `mapstructure:"Region"`
	Bucket          string `mapstructure:"Bucket"`
	Version         string `mapstructure:"Version"`
	Root            string `mapstructure:"Root"`
	Domain          string `mapstructure:"Domain"`
}

// 文件读取
func Read(filepath string) []byte {
	f, err := os.Open(filepath)
	if err != nil {
		log.Println("read file fail", err)
		return nil
	}
	defer f.Close()
	fd, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println("read to fd fail", err)
		return nil
	}

	return fd
}

// 上传本地文件
// @param savePath 保存至目录
// @param filePath 本地文件路径
// @return objectKey,domain,error
func (config Aws) UploadLoclFile(savePath, filePath string) (objectKey string, domain string, err error) {
	// 创建新的会话
	s3Config := &aws.Config{
		Region:      aws.String(config.Region), // 替换自己账户的region
		Credentials: credentials.NewStaticCredentials(config.AccessKeyId, config.SecretAccesskey, ""),
	}
	sess, err := session.NewSession(s3Config)
	if err != nil {
		return "", "", err
	}
	// 文件保存标识
	var saveDir string
	if config.Root != "" {
		saveDir += strings.Trim(config.Root, "/") + "/"
	}
	if savePath != "" {
		saveDir += strings.Trim(savePath, "/") + "/"
	}
	fileName := php2go.Uniqid("") + strings.ToLower(filepath.Ext(filePath))
	objectKey = fmt.Sprint(saveDir, fileName)
	// 读取文件流
	res := Read(filePath)
	// 上传文件
	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(config.Bucket), // bucket名称
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(res),
		ContentType: aws.String(http.DetectContentType(res)),
	})
	if err != nil {
		log.Println("PUT err", err)
		return "", "", err
	}
	return objectKey, strings.Trim(config.Domain, "/") + "/", nil
}
