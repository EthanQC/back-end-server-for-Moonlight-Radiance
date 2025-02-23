// 用户请求的处理

package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterHandler 处理用户注册请求
func RegisterHandler(c *gin.Context) {
	var input RegisterInput

	err := c.Copy().ShouldBindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalit input"})
		return
	}

	err = Register(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// LoginHandler 处理用户登录请求
func LoginHandler(c *gin.Context) {
	var input LoginInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := Login(input)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
