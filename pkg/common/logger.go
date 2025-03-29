// 日志记录模块

package common

import (
	"io"
	"log"
	"os"
)

// InitLogger 初始化日志记录器
func InitLogger() {
	file, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("could not open log file: %v", err)
	}
	// 使用 MultiWriter 将日志同时写到文件和控制台
	multiWriter := io.MultiWriter(file, os.Stdout)

	// 设置标准库 log 的输出目标为 multiWriter
	log.SetOutput(multiWriter)
	// 设置日志格式（带日期、时间、文件名和行号）
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// 设置日志前缀
	log.SetPrefix("INFO: ")
}
