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
	MaxOpenConns    int    `mapstructure:"MaxOpenConns"`    // 设置最大打开的连接数，避免连接过多导致数据库压力过大
	MaxIdleConns    int    `mapstructure:"MaxIdleConns"`    // 设置连接池中的最大闲置连接数
	ConnMaxLifetime int    `mapstructure:"ConnMaxLifetime"` // 设置连接的最大生存期，防止使用过期连接(单位秒)
	ConnMaxIdleTime int    `mapstructure:"ConnMaxIdleTime"` // 设置连接在闲置状态的最大存活时间(单位秒)
	LogLevel        int    `mapstructure:"LogLevel"`        // SQL日志级别(Silent:1, Error:2, Warn:3, Info:4)
	gormDb          *gorm.DB
}

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
		panic("failed to connect database:" + config.DataBase + fmt.Sprintf(",%s", err.Error()))
	} else if sqlDB, err := config.gormDb.DB(); err != nil {
		panic("connect db server failed.:" + config.DataBase + fmt.Sprintf(",%s", err.Error()))
	} else {
		// 数据库连接设置
		// 配置建议：
		// maxOpenConns 应根据数据库服务器的最大连接数限制和应用的并发量来设置
		// maxIdleConns 不宜设置过大，否则会占用过多数据库连接资源
		// connMaxLifetime 应小于 MySQL 配置的 wait_timeout（默认 8 小时），建议设置为几分钟
		// 生产环境中需要根据实际负载情况调整这些参数
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.ConnMaxLifetime))
		sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(config.ConnMaxIdleTime))
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
