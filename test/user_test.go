package test

import (
	"log"
	"os"
	"testing"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/common"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/user"
	"github.com/stretchr/testify/assert"
)

// TestMain 测试初始化
func TestMain1(m *testing.M) {
	// 测试前初始化数据库（一般使用测试专用的DB，如内存或测试库）
	err := common.InitDB("test_user:password@/test_db?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal("Test DB init failed")
	}
	// 运行所有测试
	code := m.Run()
	// 测试结束清理资源
	common.CloseDB()
	os.Exit(code)
}

// TestUserRegister 测试注册功能
func TestUserRegister(t *testing.T) {
	// 构造一个注册输入
	input := user.RegisterInput{
		Username: "test_user",
		Password: "123456",
	}

	err := user.Register(input)
	assert.Nil(t, err, "注册应该成功")
}

// TestUserLogin 测试登录功能
func TestUserLogin(t *testing.T) {
	input := user.LoginInput{
		Username: "test_user",
		Password: "123456",
	}

	token, err := user.Login(input)
	assert.Nil(t, err, "登录应该成功")
	assert.NotEmpty(t, token, "登录成功后应返回Token")
}
