package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin" // Gin Web 框架
	"gorm.io/gorm"
)

// RegisterHandler 处理用户注册请求
func RegisterHandler(c *gin.Context) {
	var input RegisterInput
	err := c.Copy().ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 注册用户
	err = Register(input)
	if err != nil {
		// 如果是用户名重复
		if errors.Is(err, ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
		// 其它数据库错误
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
		// 如果是 user not found 或密码错误，可以统一返回401
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrIncorrectPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
