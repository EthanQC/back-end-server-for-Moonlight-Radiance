package main

import (
	"log"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/api/http"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/configs"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/internal/auth"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/common"
)

var cfg *configs.Config

func init() {
	// 加载配置
	cfg = configs.LoadConfig()

	// 初始化日志
	common.InitLogger()

	// 初始化数据库
	if cfg.Database.DSN == "" {
		log.Fatal("DB_DSN is required")
	}
	if err := common.InitDB(cfg.Database.DSN); err != nil {
		log.Fatal("Database initialization failed: ", err)
	}

	// 初始化 Redis
	if cfg.Redis.Addr == "" {
		log.Fatal("REDIS_ADDR is required")
	}
	common.InitRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)

	// 初始化 JWT 认证
	if cfg.JWT.Secret == "" {
		log.Fatal("JWT_SECRET is required")
	}
	auth.InitJWT(cfg.JWT.Secret)
}

func main() {
	// 设置路由
	router := http.SetupRouter()

	// 启动Web服务器
	log.Println("Server is starting on :8080")
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
