// 字符串处理
package function

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// 字符串转INT类型
func StrToInt(str string) int {
	value, _ := strconv.Atoi(str)
	return value
}

func JoinStr(i any, sep string) string {
	s := fmt.Sprintf("%v", i)
	s_slice := regexp.MustCompile(`[\w.]+`).FindAllString(s, -1)
	return strings.Join(s_slice, sep)
}

// 字符串首字母大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// 字符串首字母小写
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// 驼峰空格转下划线
func CamelToSnake(str string) string {
	str = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(str, "_")                 // 非常规字符转化为 _
	snake := regexp.MustCompile("(.)([A-Z][a-z]+)").ReplaceAllString(str, "${1}_${2}")   // 拆分出连续大写
	snake = regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(snake, "${1}_${2}") // 拆分单词
	snake = strings.ToLower(snake)                                                       // 字符串转小写
	snake = strings.ReplaceAll(snake, "__", "_")                                         // 替换重复下划线
	return snake
}

// 字符串数组匹配，包含任意一项
func StrContains(str string, strArr []string) bool {
	for _, s := range strArr {
		if strings.Contains(str, s) {
			return true
		}
	}
	return false
}

// 按照规则，参数名ASCII码从小到大排序后拼接
// data 待拼接的数据
// sep 连接符
// onlyValues 是否只包含参数值，true则不包含参数名，否则参数名和参数值均有
// includeEmpty 是否包含空值，true则包含空值，否则不包含，注意此参数不影响参数名的存在
// exceptKeys 被排除的参数名，不参与排序及拼接
func JoinStringsInASCII(data map[string]string, sep string, onlyValues, includeEmpty bool, exceptKeys ...string) string {
	var list []string
	var keyList []string
	m := make(map[string]int)
	if len(exceptKeys) > 0 {
		for _, except := range exceptKeys {
			m[except] = 1
		}
	}
	for k := range data {
		if _, ok := m[k]; ok {
			continue
		}
		value := data[k]
		if !includeEmpty && value == "" {
			continue
		}
		if onlyValues {
			keyList = append(keyList, k)
		} else {
			list = append(list, fmt.Sprintf("%s=%s", k, value))
		}
	}
	if onlyValues {
		sort.Strings(keyList)
		for _, v := range keyList {
			list = append(list, data[v])
		}
	} else {
		sort.Strings(list)
	}
	return strings.Join(list, sep)
}
