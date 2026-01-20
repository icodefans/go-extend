// Mysql
package pgsql

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PgSQL配置
type PgSQL struct {
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

// 数据库连接
func (config *PgSQL) Connect() (gormDb *gorm.DB) {
	// 实现单利模式
	if config.gormDb != gormDb {
		return config.gormDb
	}
	// 数据库连接
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.HostName, config.UserName, config.PassWord, config.DataBase, config.HostPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	// 连接池配置
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期

	// 测试连接
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	// 成功返回
	config.gormDb = db
	return config.gormDb
}
