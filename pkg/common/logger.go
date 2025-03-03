// 日志记录模块

package common

import (
	"log"
	"os"
)

var logger *log.Logger

// InitLogger 初始化日志记录器
func InitLogger() {
	file, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("could not open log file: %v", err)
	}
	logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// LogMessage 记录日志
func LogMessage(message string) {
	logger.Println(message)
}
