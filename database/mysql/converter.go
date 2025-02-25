package mysql

import (
	"fmt"

	"github.com/icodefans/go-extend/vendor2/converter"
)

// mysql生成struct结构体实例
func (config *MySQL) Table2Struct(tableName string) error {
	if tableName == "" {
		return fmt.Errorf("表名不能为空")
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.UserName,
		config.PassWord,
		config.HostName,
		config.HostPort,
		config.DataBase,
	)
	savePath := fmt.Sprintf("./runtime/model/%s.go", tableName)

	// 初始化
	t2t := converter.NewTable2Struct()
	// 个性化配置
	t2t.Config(&converter.T2tConfig{
		// 如果字段首字母本来就是大写, 就不添加tag, 默认false添加, true不添加
		RmTagIfUcFirsted: false,
		// tag的字段名字是否转换为小写, 如果本身有大写字母的话, 默认false不转
		TagToLower: false,
		// 字段首字母大写的同时, 是否要把其他字母转换为小写,默认false不转换
		UcFirstOnly: true,
		// // 每个struct放入单独的文件,默认false,放入同一个文件(暂未提供)
		// SeperatFile: false,
		// 是否把字段备注设置为注释
		// ColumnComment: true,
	})
	// 开始迁移转换
	err := t2t.
		// 指定某个表,如果不指定,则默认全部表都迁移
		Table(tableName).
		// 表前缀
		Prefix("").
		// 是否添加json tag
		EnableJsonTag(true).
		// 生成struct的包名(默认为空的话, 则取名为: package model)
		PackageName("").
		// tag字段的key值,默认是orm
		TagKey("gorm").
		// 是否添加结构体方法获取表名
		RealNameMethod("").
		// 生成的结构体保存路径
		SavePath(savePath).
		// 数据库dsn
		Dsn(dsn).
		// 执行
		Run()

	return err
}
