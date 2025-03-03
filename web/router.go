package web

import (
	"net/http"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/game"
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

	// 游戏相关路由
	gameGroup := r.Group("/game")
	{
		// 简单示例：获取随机卡牌
		gameGroup.GET("/card/random", func(c *gin.Context) {
			dummyUser := &user.User{ID: 1, Username: "TestUser"}
			card := game.CardEngine(dummyUser)
			c.JSON(http.StatusOK, gin.H{
				"card": card,
			})
		})
		// 其他与游戏相关的接口
	}

	return r
}
