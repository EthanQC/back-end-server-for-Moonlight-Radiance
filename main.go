package main

import (
	"log"
	"os"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/auth"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/common"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/game"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/user"
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
	// 创建Gin路由
	router := gin.Default()

	// 用户登录与注册API
	router.POST("/register", user.RegisterHandler)
	router.POST("/login", user.LoginHandler)

	// 示例：卡牌引擎测试
	router.GET("/card/:userId", func(c *gin.Context) {
		userId := c.Param("userId")
		// 假设你会从数据库加载用户数据
		// user := getUserById(userId)
		player := &user.User{
			ID:       1, // 示例ID
			Username: "TestUser",
		}
		card := game.CardEngine(player)
		c.JSON(200, gin.H{
			"user_id": userId,
			"card":    card,
		})
	})

	// 启动Web服务器
	log.Println("Server is starting on :8080")
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
