package web

import (
	"net/http"
	"strings"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/auth"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 简单的JWT验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
			return
		}

		userID, err := auth.ParseJWT(tokenParts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// 如果需要，可以将 userID 注入到上下文中，方便后续业务处理
		c.Set("user_id", userID)
		c.Next()
	}
}
