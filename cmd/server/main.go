package main

import (
	"log"
	"net/http"
	"os"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/internal/auth"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/internal/user"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/common"
	"github.com/gin-gonic/gin"
)

func init() {
	// 初始化日志
	common.InitLogger()

	// 读取环境变量或配置文件，初始化数据库和Redis
	dbDSN := os.Getenv("DB_DSN") // MySQL连接字符串
	if dbDSN == "" {
		log.Fatal("DB_DSN is required")
	}
	err := common.InitDB(dbDSN)
	if err != nil {
		log.Fatal("Database initialization failed: ", err)
	}

	redisAddr := os.Getenv("REDIS_ADDR") // Redis服务器地址
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR is required")
	}
	common.InitRedis(redisAddr, "", 0) // 默认不需要密码，数据库0

	// 初始化JWT认证
	auth.InitJWT(os.Getenv("JWT_SECRET"))
}

func main() {
	// 设置为发布模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由器时设置信任代理
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// 用户登录与注册API
	router.POST("/register", user.RegisterHandler)
	router.POST("/login", user.LoginHandler)

	// 添加静态文件支持
	router.Static("/frontend", "./frontend")

	// 添加根路由
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/frontend/html/login.html")
	})

	// 启动Web服务器
	log.Println("Server is starting on :8080")
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
