package mysql

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/icodefans/go-extend/function"

	"github.com/syyongx/php2go"
)

type Where [][3]any

// 自定义where条件解析逻辑
func WhereParse(where Where, fields ...string) (whereSQL string, vals []any, err error) {
	if len(fields) == 0 {
		return "", nil, fmt.Errorf("where字段限制不能为空")
	}
	var (
		operator = []string{"=", "<", ">", "<>", "<=", ">=", "is", "in", "not in", "like", "find_in_set", "overlaps"}
	)
	for _, item := range where {
		var (
			field  string // 字段
			action string // 操作符
			ok     bool
			wen    = " ?"
		)
		if field, ok = item[0].(string); !ok {
			return "", nil, fmt.Errorf("where字段需要是字符串类型")
		} else if action, ok = item[1].(string); !ok {
			return "", nil, fmt.Errorf("where操作符需要是字符串类型")
		} else if field == "" {
			return "", nil, fmt.Errorf("where字段不能为空")
		} else if action == "" {
			return "", nil, fmt.Errorf("where操作符不能为空")
		} else if !function.InArrayString(action, operator) {
			return "", nil, fmt.Errorf("where不支持该操作符:%v", action)
		} else if len(fields) > 0 && !php2go.InArray(field, fields) {
			return "", nil, fmt.Errorf("where不支持该字段:%v", field)
		} else if _, ok = item[2].([]any); !(action == "in" || action == "not in" || action == "overlaps") && ok {
			return "", nil, fmt.Errorf("where字段%s操作符%s，类型必须是值类型%v", item[0], item[1], item[2])
		} else if strings.Contains(field, "->") { // json对象查询 margin->'$.select' = 'rate'
			field = fmt.Sprintf("%s", field)
		} else if fff := strings.Split(field, "."); len(fff) > 2 {
			return "", nil, fmt.Errorf("where字段点分隔不符合规范:%v", field)
		} else if len(fff) == 1 {
			field = fmt.Sprintf("`%s`", field)
		} else if len(fff) == 2 {
			field = fmt.Sprintf("`%s`.`%s`", fff[0], fff[1])
		}
		if _, ok = item[2].([]any); ok && (action == "in" || action == "not in") {
			wen = " (?)"
			field = fmt.Sprintf("%s ", field)
			vals = append(vals, item[2])
		} else if value, ok := item[2].(string); ok && (action == "in" || action == "not in") {
			wen = " (?)"
			field = fmt.Sprintf("%s ", field)
			vals = append(vals, strings.Split(value, ","))
		} else if action == "find_in_set" {
			wen = fmt.Sprintf("(?,%s)", field)
			field = ""
			vals = append(vals, item[2])
		} else if action == "overlaps" {
			wen = ""
			action = ""
			field = fmt.Sprintf(`JSON_OVERLAPS(%s,?)=1 `, field)
			if val, err := json.Marshal(item[2]); err == nil {
				vals = append(vals, val)
			} else {
				vals = append(vals, item[2])
			}
		} else if fileds := strings.Split(field, "->"); action == "search" && len(fileds) == 2 {
			wen, action = "", ""
			subField := strings.Trim(fileds[1], `'`)
			subField = strings.TrimLeft(subField, "$.")
			field = fmt.Sprintf(`JSON_VALID(%s) = 1 AND JSON_SEARCH(%s, 'one', ?, NULL, '$[*].%s') IS NOT NULL`, fileds[0], fileds[0], subField)
			vals = append(vals, item[2])
		} else {
			field = fmt.Sprintf("%s ", field)
			vals = append(vals, item[2])
		}
		if whereSQL != "" {
			whereSQL += " AND "
		}
		whereSQL += fmt.Sprint(field, strings.ToUpper(action), wen)
	}
	return whereSQL, vals, nil
}
