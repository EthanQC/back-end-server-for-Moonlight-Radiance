// 数据库连接
package common

import (
	"log" // Go标准日志库
	"time"

	"gorm.io/driver/mysql"           // GORM 的 MySQL 驱动
	"gorm.io/gorm"                   // GORM 核心库
	gormLogger "gorm.io/gorm/logger" // GORM 提供的日志模块
)

// 包级别的全局变量 DB，表示数据库连接实例
var DB *gorm.DB

// InitDB 初始化 MySQL 数据库连接
// dsn: data source name，数据源名称
// 返回 error 类型，指示初始化过程中是否发生错误
func InitDB(dsn string) error {
	var err error

	// 保证在数据库操作中，时间数据使用统一且正确的时区
	config := &gorm.Config{

		// Default 是该模块中的一个默认日志实例，包含了一些预设的日志行为，比如输出到标准输出
		// LogMode 是 gormLogger.Default 的一个方法，用于设置日志的级别
		// gormLogger.Info 是日志级别之一，表示“信息”级别
		// GORM 的日志级别包括：
		// Silent：不记录日志
		// Error：只记录错误
		// Warn：记录警告和错误
		// Info：记录详细信息，包括执行的 SQL 语句、查询时间、错误等
		Logger: gormLogger.Default.LogMode(gormLogger.Info),

		// 获取当前时间
		NowFunc: func() time.Time {
			loc, _ := time.LoadLocation("Asia/Shanghai")
			return time.Now().In(loc)
		},
	}

	DB, err = gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}

	// 获取连接池实例
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get database instance: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接的最大可复用时间

	log.Println(("Database connected successfully"))
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() {

	db, err := DB.DB()

	if err == nil {
		db.Close()
	} else {
		log.Fatalf("Failed to close database: %v", err)
	}
}
