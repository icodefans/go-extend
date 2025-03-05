// 参数验证服务
package service

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/icodefans/go-extend/function"
)

// 全局数据验证器（默认中文语言包）
func Validate(args ...any) error {
	if len(args) == 0 {
		return fmt.Errorf("验证器参数不能为空")
	}

	// 实例化验证对象
	validate := validator.New()

	// 自定义约束注册
	if err := validate.RegisterValidation("time", CheckTime); err != nil {
		return err
	} else if err := validate.RegisterValidation("date", CheckDate); err != nil {
		return err
	}

	// 语言设置，非中文则使用英文语言包
	var AcceptLanguage = "zh-cn"

	// 多语言翻译器
	var trans ut.Translator
	if AcceptLanguage == "en-us~" { // 不翻译英文
		trans, _ = ut.New(en.New()).GetTranslator("en")
		err := en_translations.RegisterDefaultTranslations(validate, trans)
		if err != nil {
			return err
		}
	} else {
		trans, _ = ut.New(zh.New()).GetTranslator("zh")
		err := zh_translations.RegisterDefaultTranslations(validate, trans)
		if err != nil {
			return err
		}
		// 注册一个函数，获取struct tag里自定义的label作为字段名
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("label")
			return name
		})
	}

	// 结构批量验证
	var validateErr error
	for i, value := range args {
		if function.IsNil(value) {
			return fmt.Errorf("参数索引%d不能为空", i)
		} else if validateErr = validate.Struct(value); validateErr != nil {
			break
		}
	}
	if validateErr == nil {
		return nil
	} else if errs, ok := validateErr.(validator.ValidationErrors); !ok || len(errs) == 0 {
		return nil
	}

	// 返回错误信息
	for _, err := range validateErr.(validator.ValidationErrors) {
		return fmt.Errorf(err.Translate(trans))
	}

	return nil
}

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
