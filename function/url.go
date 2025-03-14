// URL处理
package function

import (
	"sort"
)

// 转化数组为URL-encode 的请求字符串
func UrlEncode(haystack map[string]string) string {
	// 排序
	var keys []string
	for k := range haystack {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// 拼接
	var dataParams string
	for _, k := range keys {
		dataParams = dataParams + k + "=" + haystack[k] + "&"
	}
	// fmt.Println(dataParams)
	ff := dataParams[0 : len(dataParams)-1]
	return ff
}
