// 用户请求的处理

package user

import (
	"net/http"

	"github.com/gin-gonic/gin" // Gin Web 框架
)

// RegisterHandler 处理用户注册请求
func RegisterHandler(c *gin.Context) {
	// 声明输入结构体对象
	var input RegisterInput

	// 绑定并验证 JSON 请求体
	// 使用 c.Copy() 复制上下文，避免并发操作修改原始上下文
	// 解析请求的 JSON 体到 RegisterInput 结构体，并触发验证
	err := c.Copy().ShouldBindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 注册用户
	err = Register(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
		return
	}

	// 返回成功响应
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
