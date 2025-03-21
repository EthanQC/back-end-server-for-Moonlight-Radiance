// 数据库连接和模型定义

package common

import (
	"log" // Go标准日志库
	"time"

	"gorm.io/driver/mysql" // GORM 的 MySQL 驱动
	"gorm.io/gorm"         // GORM核心库
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

	// 设置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get database instance: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println(("Database connected successfully"))
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() {

	db, err := DB.DB()

	if err == nil {
		db.Close()
	} else {
		log.Printf("Failed to close database: %v", err)
	}
}
