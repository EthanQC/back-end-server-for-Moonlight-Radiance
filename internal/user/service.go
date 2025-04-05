package user

import (
	"errors"
	"strings"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/internal/auth"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/common"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 自定义错误
var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrIncorrectPassword = errors.New("incorrect password")
)

// Register 用户注册
func Register(input RegisterInput) error {
	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err //哈希失败
	}

	user := User{
		Username: input.Username,
		Password: string(hashedPassword),
		Role:     "player",
	}

	// 数据库写入
	result := common.DB.Create(&user)
	if result.Error != nil {
		// 如果是重复用户名(1062是MySQL主键/唯一索引冲突)
		if strings.Contains(result.Error.Error(), "1062") {
			return ErrUserAlreadyExists
		}
		return result.Error
	}

	return nil
}

// Login 用户登录
func Login(input LoginInput) (string, error) {
	var user User
	result := common.DB.Where("username = ?", input.Username).First(&user)
	if result.Error != nil {
		// 未查到记录
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", gorm.ErrRecordNotFound
		}
		// 其它数据库错误
		return "", result.Error
	}

	// 验证密码哈希
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return "", ErrIncorrectPassword
	}

	// 生成 JWT
	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}
