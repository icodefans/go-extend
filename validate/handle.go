package validate

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// 检测日期格式是否正确格式，年月日(2022-02-01)
func CheckDate(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	pattern := `^((([0-9]{3}[1-9]|[0-9]{2}[1-9][0-9]{1}|[0-9]{1}[1-9][0-9]{2}|[1-9][0-9]{3})-(((0[13578]|1[02])-(0[1-9]|[12][0-9]|3[01]))|((0[469]|11)-(0[1-9]|[12][0-9]|30))|(02-(0[1-9]|[1][0-9]|2[0-8]))))|((([0-9]{2})(0[48]|[2468][048]|[13579][26])|((0[48]|[2468][048]|[3579][26])00))-02-29))$`
	result, _ := regexp.MatchString(pattern, value)
	return result
}

// 检测时间格式是否正确格式，时分秒(8:03:02)
func CheckTime(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	pattern := `^([0-1]?[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])$`
	result, _ := regexp.MatchString(pattern, value)
	return result
}
