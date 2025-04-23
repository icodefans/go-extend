package function

import (
	"math"
	"strconv"
	"strings"
)

// FormatCurrency 函数用于将货币金额格式化为千分位形式
func FormatCurrency(amount float64, currencySymbol string, precision int) string {
	// 将金额转换为字符串，并保留指定的小数位数
	str := strconv.FormatFloat(amount, 'f', precision, 64)
	// 分割整数部分和小数部分
	parts := strings.Split(str, ".")
	integerPart := parts[0]
	var decimalPart string
	if len(parts) > 1 {
		decimalPart = parts[1]
	}

	// 处理整数部分的千分位
	n := len(integerPart)
	if n <= 3 {
		result := currencySymbol + integerPart
		if decimalPart != "" {
			result += "." + decimalPart
		}
		return result
	}
	// 计算需要插入逗号的次数
	commas := int(math.Floor(float64(n-1) / 3))
	result := make([]rune, len(integerPart)+commas+len(currencySymbol))
	copy(result[:len(currencySymbol)], []rune(currencySymbol))
	j := len(result) - 1
	for i := n - 1; i >= 0; i-- {
		result[j] = rune(integerPart[i])
		j--
		if (n-i)%3 == 0 && i > 0 {
			result[j] = ','
			j--
		}
	}

	// 组合整数部分、小数部分和货币符号
	if decimalPart != "" {
		return string(result) + "." + decimalPart
	}
	return string(result)
}
