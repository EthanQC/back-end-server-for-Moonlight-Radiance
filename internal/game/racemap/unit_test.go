package racemap

import (
	"testing"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/common"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) func() {
	err := common.InitDB("root:your_password@tcp(localhost:3306)/moonlight_test?charset=utf8mb4&parseTime=True&loc=Local")
	assert.NoError(t, err)

	common.DB.Exec(`
        CREATE TABLE IF NOT EXISTS race_maps (
            id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
            length INT NOT NULL,
            created_at TIMESTAMP,
            updated_at TIMESTAMP
        )
    `)

	common.DB.Exec(`
        CREATE TABLE IF NOT EXISTS positions (
            race_map_id INT UNSIGNED NOT NULL,
            player_id INT UNSIGNED NOT NULL,
            location INT NOT NULL DEFAULT 0,
            moon_value FLOAT NOT NULL DEFAULT 0,
            PRIMARY KEY (race_map_id, player_id)
        )
    `)

	return func() {
		common.DB.Exec("DROP TABLE race_maps")
		common.DB.Exec("DROP TABLE positions")
	}
}

func TestRaceMapFlow(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	s := NewRaceMapService()
	playerID := uint(1)

	// 1. 测试创建地图
	t.Run("Create Map", func(t *testing.T) {
		raceMap, err := s.CreateMap()
		assert.NoError(t, err)
		assert.Equal(t, DefaultMapLength, raceMap.Length)
	})

	// 2. 测试初始化玩家位置
	t.Run("Init Player Position", func(t *testing.T) {
		raceMap, _ := s.CreateMap()

		err := s.InitPlayerPosition(raceMap.ID, playerID)
		assert.NoError(t, err)

		position, err := s.GetPosition(raceMap.ID, playerID)
		assert.NoError(t, err)
		assert.Equal(t, 0, position.Location)
		assert.Equal(t, float64(0), position.MoonValue)
	})

	// 3. 测试移动
	t.Run("Move Forward", func(t *testing.T) {
		raceMap, _ := s.CreateMap()
		_ = s.InitPlayerPosition(raceMap.ID, playerID)

		result, err := s.MoveForward(raceMap.ID, playerID, 1.5)
		assert.NoError(t, err)
		assert.Equal(t, 1, result.NewLocation)
		assert.Equal(t, float64(1.5), result.MoonValue)
		assert.False(t, result.IsFinished)

		// 测试到达终点
		for i := 0; i < DefaultMapLength-2; i++ {
			_, _ = s.MoveForward(raceMap.ID, playerID, 1.5)
		}

		finalResult, err := s.MoveForward(raceMap.ID, playerID, 1.5)
		assert.NoError(t, err)
		assert.True(t, finalResult.IsFinished)
	})
}
