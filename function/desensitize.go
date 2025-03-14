package function

import (
	"fmt"
	"strings"
)

// 字符串脱敏，默认使用*号代替
// var str = "18627092724" // 原始字符串
// var beginLen = 3        // 开头字符长度
// var endLen = 4          // 结尾字符长度
// var starLen = 4         // 脱敏字符长度
func Desensitize(str string, beginLen, endLen, starLen uint32) string {
	var starStr = strings.Repeat("*", int(starLen)) // 脱敏字符数量
	var strArr = strings.Split(str, "")             // 字符数组
	var strLen = uint32(len(strArr))                // 字符串长度
	// 默认值设置
	if beginLen > strLen { // 开头字符长度大于字符串长度，则使用字符串长度
		beginLen = strLen
	}
	if beginLen >= strLen { // 开头字符长度大于等于字符串长度，结尾字符串不显示
		endLen = 0
	} else if beginLen+endLen > strLen { // 开头加上结尾字符串大于字符串长度
		endLen = strLen - beginLen
	}
	// 字符串设置
	var beginStr = strings.Join(strArr[:beginLen], "")    // 开头字符串
	var endStr = strings.Join(strArr[strLen-endLen:], "") // 结尾字符串
	var rsStr = fmt.Sprintf("%s%s%s", beginStr, starStr, endStr)
	return rsStr
}
