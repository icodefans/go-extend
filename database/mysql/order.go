package mysql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/icodefans/go-extend/function"

	"github.com/syyongx/php2go"
)

type Order [][2]string

// 自定义order条件解析逻辑
func OrderParse(order Order, fields ...string) (orderSQL string, err error) {
	for _, item := range order {
		if item[0] == "" {
			return "", fmt.Errorf("order字段不能为空")
		} else if item[1] == "" {
			return "", fmt.Errorf("order操作符不能为空")
		} else if !php2go.InArray(item[1], []string{"asc", "desc"}) {
			return "", fmt.Errorf("order操作符不在选择范围内[asc,desc]")
		} else if len(fields) > 0 && !php2go.InArray(item[0], fields) {
			return "", fmt.Errorf("order不支持该字段:%v", item[0])
		} else if field := strings.Split(item[0], "."); !strings.Contains(item[0], "->") && len(field) > 2 {
			return "", fmt.Errorf("字段点分隔不符合规范:%v", item[0])
		} else if !strings.Contains(item[0], "->") && len(field) == 1 {
			item[0] = fmt.Sprintf("`%s`", item[0])
		} else if !strings.Contains(item[0], "->") && len(field) == 2 {
			item[0] = fmt.Sprintf("`%s`.`%s`", field[0], field[1])
		}
		if orderSQL != "" {
			orderSQL += " , "
		}
		if !function.InArrayString(item[1], []string{"asc", "desc"}) {
			return "", errors.New("order排序模式错误")
		}
		orderSQL += fmt.Sprint("", item[0], "", " ", strings.ToUpper(item[1]))
	}
	return orderSQL, nil
}
