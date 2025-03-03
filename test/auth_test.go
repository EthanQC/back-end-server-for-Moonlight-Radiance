package test

import (
	"os"
	"testing"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/auth"
	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	// 设置一个测试用的JWT_SECRET
	auth.InitJWT("test_secret")

	// 生成token
	token, err := auth.GenerateJWT(123)
	assert.Nil(t, err, "生成token不应该出错")
	assert.NotEmpty(t, token, "token不应该为空")

	// 解析token
	userID, err := auth.ParseJWT(token)
	assert.Nil(t, err, "解析token不应该出错")
	assert.Equal(t, uint(123), userID, "解析后的userID应该匹配")
}

func TestMain(m *testing.M) {
	// 可以在这里做一些Auth模块需要的初始化或Mock
	code := m.Run()
	os.Exit(code)
}
