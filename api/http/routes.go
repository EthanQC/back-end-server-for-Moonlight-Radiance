package http

import (
	"net/http"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/internal/user"
	"github.com/gin-gonic/gin"
)

// SetupRouter 设置并返回 Gin 引擎
func SetupRouter() *gin.Engine {

	// 创建路由器时设置信任代理
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// 静态文件支持
	router.Static("/frontend", "./frontend")

	// 根路由组
	root := router.Group("/")
	{
		// 健康检查
		root.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
		// 重定向到登录页
		root.GET("", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/frontend/html/login.html")
		})

		// 用户登录和注册
		root.POST("/login", user.LoginHandler)
		root.POST("/register", user.RegisterHandler)

	}

	return router
}
