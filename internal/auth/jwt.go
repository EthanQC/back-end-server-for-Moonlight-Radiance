// JWT认证的生成与验证
// JWT（JSON Web Token）：一种用于身份验证的令牌格式，通常包含用户ID、过期时间等

package auth

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret []byte

// InitJWT 初始化JWT密钥
func InitJWT(secret string) {
	jwtSecret = []byte(secret)
}

// GenerateJWT 生成JWT token
func GenerateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24小时过期
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)

	if err != nil {
		log.Printf("Error generating token: %v", err)
		return "", err
	}

	return signedToken, nil
}

// ParseJWT 解析JWT token
func ParseJWT(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return 0, jwt.ErrSignatureInvalid
	}

	userID, ok := claims["user_id"].(float64)

	if !ok {
		return 0, jwt.ErrSignatureInvalid
	}

	return uint(userID), nil
}
