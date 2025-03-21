package room

import (
	"testing"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/common"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) func() {
	err := common.InitDB("root:your_password@tcp(localhost:3306)/moonlight_test?charset=utf8mb4&parseTime=True&loc=Local")
	assert.NoError(t, err)

	// 创建测试表
	common.DB.Exec(`
        CREATE TABLE IF NOT EXISTS rooms (
            id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
            status TINYINT NOT NULL,
            player1_id INT UNSIGNED NOT NULL,
            player2_id INT UNSIGNED,
            current_player_id INT UNSIGNED NOT NULL,
            battle_map_id INT UNSIGNED NOT NULL,
            race_map_id INT UNSIGNED NOT NULL,
            created_at TIMESTAMP,
            updated_at TIMESTAMP
        )
    `)

	common.DB.Exec(`
        CREATE TABLE IF NOT EXISTS player_progress (
            room_id INT UNSIGNED NOT NULL,
            player_id INT UNSIGNED NOT NULL,
            position INT NOT NULL DEFAULT 0,
            moon_value FLOAT NOT NULL DEFAULT 0,
            PRIMARY KEY (room_id, player_id)
        )
    `)

	return func() {
		common.DB.Exec("DROP TABLE rooms")
		common.DB.Exec("DROP TABLE player_progress")
	}
}

func TestRoomFlow(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	s := NewRoomService()
	player1ID := uint(1)
	player2ID := uint(2)

	// 1. 测试创建房间
	t.Run("Create Room", func(t *testing.T) {
		room, err := s.CreateRoom(player1ID)
		assert.NoError(t, err)
		assert.Equal(t, player1ID, room.Player1ID)
		assert.Equal(t, Waiting, room.Status)

		progress, err := s.GetPlayerProgress(room.ID, player1ID)
		assert.NoError(t, err)
		assert.Equal(t, 0, progress.Position)
	})

	// 2. 测试加入房间
	t.Run("Join Room", func(t *testing.T) {
		room, _ := s.CreateRoom(player1ID)

		err := s.JoinRoom(room.ID, player2ID)
		assert.NoError(t, err)

		updatedRoom, _ := s.GetRoomState(room.ID)
		assert.Equal(t, Playing, updatedRoom.Status)
		assert.Equal(t, &player2ID, updatedRoom.Player2ID)

		progress, err := s.GetPlayerProgress(room.ID, player2ID)
		assert.NoError(t, err)
		assert.Equal(t, 0, progress.Position)
	})

	// 3. 测试回合结束
	t.Run("End Turn", func(t *testing.T) {
		room, _ := s.CreateRoom(player1ID)
		_ = s.JoinRoom(room.ID, player2ID)

		err := s.EndTurn(room.ID, player1ID)
		assert.NoError(t, err)

		updatedRoom, _ := s.GetRoomState(room.ID)
		assert.Equal(t, player2ID, updatedRoom.CurrentPlayerID)
	})
}
