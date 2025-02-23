// 数据库连接和模型定义

package common

import (
	"log" // Go标准日志库

	"gorm.io/driver/mysql" // GORM 的 MySQL 驱动
	"gorm.io/gorm"         // GORM核心库
)

var DB *gorm.DB

// InitDB 初始化 MySQL 数据库连接
// dsn: data source name，数据源名称
func InitDB(dsn string) error {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}
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
