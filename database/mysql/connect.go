// Mysql
package mysql

import (
	// myLogger "github.com/icodefans/go-extend/logger"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQL配置
type MySQL struct {
	HostName        string `mapstructure:"HostName"`        // 服务器地址
	HostPort        string `mapstructure:"HostPort"`        // 端口
	DataBase        string `mapstructure:"DataBase"`        // 数据库名
	UserName        string `mapstructure:"UserName"`        // 用户名
	PassWord        string `mapstructure:"PassWord"`        // 密码
	Charset         string `mapstructure:"Charset"`         // 数据库编码默认采用utf8
	Prefix          string `mapstructure:"Prefix"`          // 数据库表前缀
	Timeout         int    `mapstructure:"Timeout"`         // 设置最大打开的连接数，默认值为0表示不限制
	MaxOpenConns    int    `mapstructure:"MaxOpenConns"`    // 设置最大打开的连接数，默认值为0表示不限制
	MaxIdleConns    int    `mapstructure:"MaxIdleConns"`    // 连接池里面允许Idel的最大连接数
	ConnMaxLifetime int    `mapstructure:"ConnMaxLifetime"` // 设置一个连接的最长生命周期
	LogLevel        int    `mapstructure:"LogLevel"`        // SQL日志级别(Silent:1, Error:2, Warn:3, Info:4)
	gormDb          *gorm.DB
}

// SetMaxOpenConns: 设置最大打开的连接数，默认值为0表示不限制。控制应用于数据库建立连接的数量，避免过多连接压垮数据库。
// SetMaxIdleConns: 连接池里面允许Idel(空闲)的最大连接数, 这些Idel的连接 就是并发时可以同时获取的连接,也是用完后放回池里面的互用的连接, 从而提升性能。
// SetConnMaxLifetime: 设置一个连接的最长生命周期，因为数据库本身对连接有一个超时时间的设置，如果超时时间到了数据库会单方面断掉连接，此时再用连接池内的连接进行访问就会出错, 因此这个值往往要小于数据库本身的连接超时时间
func (config *MySQL) Connect() (gormDb *gorm.DB) {
	// 实现单利模式
	if config.gormDb != gormDb {
		return config.gormDb
	}
	// 数据库连接
	var err error
	if config.gormDb, err = gorm.Open(mysql.Open(fmt.Sprintf( // 连接配置
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=%dms",
		config.UserName, config.PassWord, config.HostName, config.HostPort,
		config.DataBase, config.Charset, config.Timeout,
	)), &gorm.Config{ // GORM配置
		SkipDefaultTransaction: true,  // 禁用默认事务
		PrepareStmt:            false, // 缓存 Prepared Statement
		Logger: logger.New( // 日志配置
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			// myLogger.NewMyWriter(),
			logger.Config{
				SlowThreshold:             time.Second * 3,                  // 慢SQL阈值(3秒)
				LogLevel:                  logger.LogLevel(config.LogLevel), // #SQL日志级别(Silent:1, Error:2, Warn:3, Info:4)
				IgnoreRecordNotFoundError: true,                             // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  false,                            // 禁用彩色打印
			},
		),
	}); err != nil {
		panic("failed to connect database:" + config.DataBase)
	} else if sqlDB, err := config.gormDb.DB(); err != nil {
		panic("connect db server failed.:" + config.DataBase)
	} else { // 数据库连接设置
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.ConnMaxLifetime))
	}
	// 成功返回
	return config.gormDb
}

// 设置最大连接数
func (config *MySQL) SetMaxOpenConns(n int) error {
	if config.gormDb == nil {
		return fmt.Errorf("数据库还未连接")
	}
	sqlDB, err := config.gormDb.DB()
	if err != nil {
		return err
	}
	// 数据库连接配置
	sqlDB.SetMaxOpenConns(n)
	return nil
}
