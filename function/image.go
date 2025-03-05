package function

import (
	"bytes"
	"encoding/hex"
	"errors"
	"os"
	"strconv"
	"strings"
)

// 先从文件流 []byte中拷贝一份前20字节进行判断
var picMap = map[string]string{
	"ffd8ffe0": "jpg",
	"ffd8ffe1": "jpg",
	"ffd8ffe8": "jpg",
	"89504e47": "png",
}

// 根据文件路径验证是否文件类型
func ImageTypeCheck(filePath string) (ok bool, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, errors.New("文件打开错误")
	}
	defer file.Close()
	result := judgeType(file)
	return result, nil
}

// 先从文件流 []byte中拷贝一份前20字节进行判断
func judgeType(file *os.File) bool {
	buf := make([]byte, 20)
	n, _ := file.Read(buf)

	fileCode := bytesToHexString(buf[:n])
	for k, _ := range picMap {
		if strings.HasPrefix(fileCode, k) {
			return true
		}
	}
	return false
}

// 获取16进制
func bytesToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	i, length := 100, len(src)
	if length < i {
		i = length
	}
	for j := 0; j < i; j++ {
		sub := src[j] & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}
