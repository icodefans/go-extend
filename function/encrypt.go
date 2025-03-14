// 加密函数
package function

import (
	"crypto/md5"
	"fmt"
)

// MD5加密
func Md5Encode(str string) string {
	// 字符串转byte
	data := []byte(str)
	// 将编码转换为字符串
	return fmt.Sprintf("%x", md5.Sum(data))
}
