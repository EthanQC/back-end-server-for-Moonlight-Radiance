package battlemap

import (
	"testing"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/common"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) func() {
	// 设置测试数据库
	err := common.InitDB("root:your_password@tcp(localhost:3306)/moonlight_test?charset=utf8mb4&parseTime=True&loc=Local")
	assert.NoError(t, err)

	// 创建测试表
	common.DB.Exec(`
        CREATE TABLE IF NOT EXISTS battle_maps (
            id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
            game_id INT UNSIGNED NOT NULL,
            size INT NOT NULL,
            grids JSON,
            created_at BIGINT,
            updated_at BIGINT
        )
    `)

	// 返回清理函数
	return func() {
		common.DB.Exec("DROP TABLE battle_maps")
	}
}

func TestBattleMapFlow(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	s := NewBattleMapService()
	gameID := uint(1)
	player1ID := uint(1)
	player2ID := uint(2)

	// 1. 测试创建地图
	t.Run("Create Map", func(t *testing.T) {
		battleMap, err := s.CreateMap(gameID)
		assert.NoError(t, err)
		assert.Equal(t, gameID, battleMap.GameID)
		assert.Equal(t, StandardMapSize, battleMap.Size)

		state, err := s.GetMapState(battleMap.ID)
		assert.NoError(t, err)
		assert.Len(t, state.Grids, StandardMapSize*StandardMapSize)
	})

	// 2. 测试放置卡牌
	t.Run("Place Card", func(t *testing.T) {
		battleMap, _ := s.CreateMap(gameID)

		// 初始化玩家的卡牌状态
		err := s.cardService.InitializePlayerDeck(gameID, player1ID)
		assert.NoError(t, err)
		err = s.cardService.DrawInitialCards(gameID, player1ID)
		assert.NoError(t, err)

		// 获取玩家手牌
		state, err := s.cardService.GetCardState(gameID, player1ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, state.Self.HandCards)

		// 尝试放置第一张手牌
		cardID := state.Self.HandCards[0].ID
		err = s.PlaceCard(battleMap.ID, player1ID, 0, 0, cardID)
		assert.NoError(t, err)

		mapState, err := s.GetMapState(battleMap.ID)
		assert.NoError(t, err)
		assert.NotNil(t, mapState.Grids[0].Placement)
		assert.Equal(t, cardID, mapState.Grids[0].Placement.CardID)
		assert.Equal(t, player1ID, mapState.Grids[0].Placement.PlayerID)

		// 测试放置到已占用格子
		err = s.PlaceCard(battleMap.ID, player2ID, 0, 0, cardID)
		assert.Error(t, err)
	})
}
