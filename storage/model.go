package storage

// 存储接口定义
type Storage interface {
	UploadLoclFile(savePath, filePath string) (objectKey string, domain string, err error)
	DownloadLocalFile(objectKey, tempFile string) (err error)
}
