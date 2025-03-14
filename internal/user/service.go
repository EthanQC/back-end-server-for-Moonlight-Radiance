// 用户服务的核心逻辑

package user

import (
	"errors" //标准错误处理

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/auth"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/common"
	"golang.org/x/crypto/bcrypt" //密码哈希库
)

// Register 用户注册
func Register(input RegisterInput) error {

	// 密码哈希处理，将明文密码转换为不可逆的 Bcrypy 哈希
	// []byte(input.Password)：将字符串密码转为字节切片
	// bcrypt.DefaultCost：哈希计算强度（默认值 10，值越大越安全但越慢）
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return err //哈希失败（如密码过长）
	}

	// 构建用户对象
	user := User{
		Username: input.Username,
		Password: string(hashedPassword), // 将返回的 []byte 类型的哈希密码转换为 string 类型
		Role:     "player",               // 默认角色是玩家
	}

	// 数据库写入
	// & 传递指针确保 GORM 可修改对象（如填充ID）
	result := common.DB.Create(&user)

	if result.Error != nil {
		return result.Error // 数据库错误（如用户名重复）
	}

	return nil
}

// Login 用户登录
func Login(input LoginInput) (string, error) {
	var user User

	// 按用户名查询用户
	result := common.DB.Where("username = ?", input.Username).First(&user)

	if result.Error != nil {
		return "", errors.New("user not found")
	}

	// 验证密码哈希
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	if err != nil {
		return "", errors.New("incorrect password")
	}

	// 生成 JWT 令牌
	token, err := auth.GenerateJWT(user.ID)

	if err != nil {
		return "", err // JWT 生成失败（如密钥错误）
	}

	return token, nil
}
