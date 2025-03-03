package test

import (
	"testing"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/game"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/user"
	"github.com/stretchr/testify/assert"
)

// TestCardEngine 测试卡牌引擎的随机性或逻辑
func TestCardEngine(t *testing.T) {
	dummyUser := &user.User{
		ID:       1,
		Username: "DummyPlayer",
	}
	card := game.CardEngine(dummyUser)
	assert.NotEmpty(t, card, "卡牌返回值不应该为空")
}

// TestPVP 测试PVP逻辑
func TestPVP(t *testing.T) {
	player1 := &user.User{ID: 1, Username: "Player1"}
	player2 := &user.User{ID: 2, Username: "Player2"}

	result, err := game.PVP(player1, player2)
	assert.Nil(t, err, "PVP对战不应该出错")
	assert.NotEmpty(t, result, "PVP结果应该有返回值")
}
