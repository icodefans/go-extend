package mysql

import (
	"fmt"
	"strings"

	"github.com/icodefans/go-extend/function"
)

type Extend [][3]string

// 自定义where条件解析逻辑
func ExtendParse(where Extend, extendField string) (whereSQL string, vals []any, err error) {
	var (
		operator = []string{"=", "<", ">", "<>", "<=", ">=", "in", "not in", "like", "find_in_set"}
	)
	for _, item := range where {
		if item[0] == "" {
			return "", nil, fmt.Errorf("where字段不能为空")
		} else if item[1] == "" {
			return "", nil, fmt.Errorf("where操作符不能为空")
		} else if !function.InArrayString(item[1], operator) {
			return "", nil, fmt.Errorf("where不支持该操作符:%v", item[1])
		}
		item[0] = fmt.Sprintf("%s->'$.%s'", extendField, item[0])
		if whereSQL != "" {
			whereSQL += " AND "
		}
		wen := " ?"
		if item[1] == "in" || item[1] == "not in" {
			wen = " (?)"
			item[0] = fmt.Sprint("", item[0], " ")
			vals = append(vals, strings.Split(item[2], ","))
		} else if item[1] == "find_in_set" {
			wen = fmt.Sprint("(?,", item[0], ")")
			item[0] = ""
			vals = append(vals, item[2])
		} else {
			item[0] = fmt.Sprint("", item[0], " ")
			vals = append(vals, item[2])
		}
		whereSQL += fmt.Sprint(item[0], strings.ToUpper(item[1]), wen)
	}
	return whereSQL, vals, nil
}
