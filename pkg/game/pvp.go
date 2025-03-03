package game

import (
	"errors"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/user"
)

// PVP 玩家对战
func PVP(player1, player2 *user.User) (string, error) {
	if player1.ID == player2.ID {
		return "", errors.New("players must be different")
	}

	// 假设卡牌机制或其他逻辑决定对战胜负
	player1Card := CardEngine(player1)
	player2Card := CardEngine(player2)

	if player1Card == player2Card {
		return "Draw", nil
	}

	// 简单的胜负判断逻辑
	if player1Card == "Sun" && player2Card != "Moon" {
		return "Player 1 wins", nil
	}

	return "Player 2 wins", nil
}
