package storage

// 引入依赖包
import (
	"fmt"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/icodefans/go-extend/function"
	"github.com/syyongx/php2go"
)

// 华为云obs配置
type HuaWei struct {
	AccessKeyId     string `mapstructure:"AccessKeyId"`
	AccessKeySecret string `mapstructure:"AccessKeySecret"`
	EndPoint        string `mapstructure:"EndPoint"`
	BucketName      string `mapstructure:"BucketName"`
	StorageRoot     string `mapstructure:"StorageRoot"`
	StorageDomain   string `mapstructure:"StorageDomain"`
	ObsClient       *obs.ObsClient
	Error           error
}

// 实例化客户端
func (huawei *HuaWei) Client() *HuaWei {
	huawei.ObsClient, huawei.Error = obs.New(huawei.AccessKeyId, huawei.AccessKeySecret, huawei.EndPoint)
	return huawei
}

// 上传本地文件
// @param root string 保存至目录
// @param filePath 本地文件路径
// @return string,error
func (huawei HuaWei) UploadLoclFile(savePath, filePath string) (objectKey string, domain string, err error) {

	// 依次填写Object的完整路径（例如exampledir/exampleobject.txt）和本地文件的完整路径（例如D:\\localpath\\examplefile.txt）。
	var saveDir string
	if huawei.StorageRoot != "" {
		saveDir += strings.Trim(huawei.StorageRoot, "/") + "/"
	}
	if savePath != "" {
		saveDir += strings.Trim(savePath, "/") + "/"
	}
	fileName := php2go.Uniqid("") + strings.ToLower(function.FileGetExt(filePath))
	objectKey = fmt.Sprint(saveDir, fileName)

	// 上传文件
	input := &obs.PutFileInput{}
	input.Bucket = huawei.BucketName
	input.Key = objectKey
	input.SourceFile = filePath
	_, err = huawei.ObsClient.PutFile(input)

	if err != nil {
		fmt.Println("Error:", err)
		return "", "", err
	}

	// 成功返回
	return objectKey, strings.Trim(huawei.StorageDomain, "/") + "/", nil
}

// 下载文件至指定位置
// @param objectKey string 对象文件
// @param tempFile 临时文件
// @return error
func (huawei HuaWei) DownloadLocalFile(objectKey, tempFile string) (err error) {
	// 参数设置
	input := &obs.DownloadFileInput{}
	input.Bucket = huawei.BucketName
	input.Key = objectKey
	input.DownloadFile = tempFile    // localfile is the full path to which objects are downloaded.
	input.EnableCheckpoint = true    // Enable the resumable download mode.
	input.PartSize = 9 * 1024 * 1024 // Set the part size to 9 MB.
	input.TaskNum = 5                // Specify the maximum number of parts that can be concurrently downloaded.

	// 文件下载
	output, err := huawei.ObsClient.DownloadFile(input)

	// 调试输出
	if obsError, ok := err.(obs.ObsError); ok {
		fmt.Printf("Code:%s\n", obsError.Code)
		fmt.Printf("Message:%s\n", obsError.Message)
	} else if err != nil {
		fmt.Println(objectKey)
		fmt.Println(err.Error())
	} else {
		fmt.Printf("RequestId:%s\n", output.RequestId)
	}

	// 成功返回
	return nil
}
