package web

import (
	"net/http"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/user"
	"github.com/gin-gonic/gin"
)

// SetupRouter 设置并返回 Gin 引擎
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 用户相关路由
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", user.RegisterHandler)
		userGroup.POST("/login", user.LoginHandler)
		// 其他与用户相关的接口
	}

	return r
}
