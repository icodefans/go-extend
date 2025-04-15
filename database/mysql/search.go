package mysql

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/syyongx/php2go"
	"github.com/icodefans/go-extend/function"
)

// 搜索结构
type Search struct {
	Where  [][3]any
	Extend [][3]string
	Order  [][2]string
}

// 自定义where条件解析逻辑
func (search *Search) WhereParse(fields ...string) (whereSQL string, vals []any, err error) {
	if len(fields) == 0 {
		return "", nil, fmt.Errorf("where字段限制不能为空")
	}
	var (
		operator = []string{"=", "<", ">", "<>", "<=", ">=", "is", "in", "not in", "like", "find_in_set", "overlaps"}
	)
	for _, item := range search.Where {
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

// 自定义where条件解析逻辑
func (search *Search) ExtendParse(extendField string) (whereSQL string, vals []any, err error) {
	var (
		operator = []string{"=", "<", ">", "<>", "<=", ">=", "in", "not in", "like", "find_in_set"}
	)
	for _, item := range search.Extend {
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

// 自定义order条件解析逻辑
func (search *Search) OrderParse(fields ...string) (orderSQL string, err error) {
	for _, item := range search.Order {
		if item[0] == "" {
			return "", fmt.Errorf("order字段不能为空")
		} else if item[1] == "" {
			return "", fmt.Errorf("order操作符不能为空")
		} else if !php2go.InArray(item[1], []string{"asc", "desc"}) {
			return "", fmt.Errorf("order操作符不在选择范围内[asc,desc]")
		} else if len(fields) > 0 && !php2go.InArray(item[0], fields) {
			return "", fmt.Errorf("order不支持该字段:%v", item[0])
		} else if field := strings.Split(item[0], "."); len(field) > 2 {
			return "", fmt.Errorf("字段点分隔不符合规范:%v", item[0])
		} else if len(field) == 1 {
			item[0] = fmt.Sprintf("`%s`", item[0])
		} else if len(field) == 2 {
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

// 自定义like条件解析逻辑
func (search *Search) LikeParse(keyword string, fields ...string) (likeSQL string) {
	if keyword == "" || len(fields) == 0 {
		return ""
	}
	likeFields := []string{}
	for _, name := range fields {
		likeFields = append(likeFields, fmt.Sprintf("%s like '%%%s%%'", name, keyword))
	}
	return strings.Join(likeFields, " OR ")
}

