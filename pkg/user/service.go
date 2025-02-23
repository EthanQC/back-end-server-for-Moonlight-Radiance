// 用户服务的核心逻辑

package user

import (
	"errors"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/auth"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/common"
	"golang.org/x/crypto/bcrypt"
)

// Register 用户注册
func Register(input RegisterInput) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user := User{
		Username: input.Username,
		Password: string(hashedPassword),
		Role:     "player", // 默认角色是玩家
	}

	result := common.DB.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Login 用户登录
func Login(input LoginInput) (string, error) {
	var user User
	result := common.DB.Where("username = ?", input.Username).First(&user)

	if result.Error != nil {
		return "", errors.New("user not found")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	if err != nil {
		return "", errors.New("incorrect password")
	}

	token, err := auth.GenerateJWT(user.ID)

	if err != nil {
		return "", err
	}

	return token, nil
}
