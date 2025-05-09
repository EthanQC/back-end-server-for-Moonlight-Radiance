package main

import (
	"log"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/api/http"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/configs"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/internal/auth"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/common"
	"github.com/gin-gonic/gin"
)

var cfg *configs.Config

func init() {
	// 加载配置
	cfg = configs.LoadConfig()
	gin.SetMode(gin.ReleaseMode)

	// 初始化日志
	common.InitLogger()

	// 初始化数据库
	if cfg.Database.DSN == "" {
		log.Fatalf("DB_DSN is required")
	}
	if err := common.InitDB(cfg.Database.DSN); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// 初始化 Redis
	if cfg.Redis.Addr == "" {
		log.Fatalf("REDIS_ADDR is required")
	}
	common.InitRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)

	// 初始化 JWT 认证
	if cfg.JWT.Secret == "" {
		log.Fatalf("JWT_SECRET is required")
	}
	auth.InitJWT(cfg.JWT.Secret)
}

func main() {
	// 设置路由
	router := http.SetupRouter()

	// 启动Web服务器
	log.Println("Server is starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
