package battlemap

import (
	"encoding/json"
	"errors"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/internal/card"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/common"
	"gorm.io/gorm"
)

const (
	StandardMapSize = 8 // 标准地图大小8x8
)

type BattleMapService struct {
	db          *gorm.DB
	cardService *card.CardService
}

func NewBattleMapService() *BattleMapService {
	return &BattleMapService{
		db:          common.DB,
		cardService: card.NewCardService(),
	}
}

// CreateMap 创建新对战地图(固定8x8)
func (s *BattleMapService) CreateMap(gameID uint) (*BattleMap, error) {
	// 初始化格子
	grids := make([]Grid, StandardMapSize*StandardMapSize)
	for i := range grids {
		grids[i] = Grid{
			Index: i,
			Position: Position{
				X: i % StandardMapSize,
				Y: i / StandardMapSize,
			},
		}
	}

	gridsJSON, err := json.Marshal(grids)
	if err != nil {
		return nil, err
	}

	battleMap := &BattleMap{
		GameID: gameID,
		Size:   StandardMapSize,
		Grids:  gridsJSON,
	}

	if err := s.db.Create(battleMap).Error; err != nil {
		return nil, err
	}

	return battleMap, nil
}

// PlaceCard 在格子上放置卡牌
func (s *BattleMapService) PlaceCard(mapID uint, playerID uint, x int, y int, cardID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 检查地图是否存在
		var battleMap BattleMap
		if err := tx.First(&battleMap, mapID).Error; err != nil {
			return err
		}

		// 检查位置是否有效
		if x < 0 || x >= StandardMapSize || y < 0 || y >= StandardMapSize {
			return errors.New("invalid grid position")
		}

		var grids []Grid
		if err := json.Unmarshal(battleMap.Grids, &grids); err != nil {
			return err
		}

		idx := y*StandardMapSize + x
		grid := &grids[idx]

		if grid.Placement != nil {
			return errors.New("grid already occupied")
		}

		// 尝试打出这张牌
		if err := s.cardService.PlayCard(battleMap.GameID, playerID, cardID); err != nil {
			return err
		}

		// 放置卡牌到地图上
		grid.Placement = &CardPlacement{
			CardID:    cardID,
			PlayerID:  playerID,
			MoonValue: 1.0, // 这个值后续需要根据游戏规则计算
		}

		gridsJSON, err := json.Marshal(grids)
		if err != nil {
			return err
		}

		return tx.Model(&battleMap).Update("grids", gridsJSON).Error
	})
}

// GetMapState 获取地图状态
func (s *BattleMapService) GetMapState(mapID uint) (*MapState, error) {
	var battleMap BattleMap
	if err := s.db.First(&battleMap, mapID).Error; err != nil {
		return nil, err
	}

	var grids []Grid
	if err := json.Unmarshal(battleMap.Grids, &grids); err != nil {
		return nil, err
	}

	// 构建玩家状态map
	playerStates := make(map[uint]*PlayerState)

	// 计算每个玩家的状态
	for _, grid := range grids {
		if grid.Placement != nil {
			if _, exists := playerStates[grid.Placement.PlayerID]; !exists {
				playerStates[grid.Placement.PlayerID] = &PlayerState{
					ID:        grid.Placement.PlayerID,
					MoonValue: 0,
				}
			}
			playerStates[grid.Placement.PlayerID].MoonValue += grid.Placement.MoonValue
		}
	}

	// 转换为数组
	players := make([]PlayerState, 0, len(playerStates))
	for _, state := range playerStates {
		players = append(players, *state)
	}

	return &MapState{
		MapID:   mapID,
		Grids:   grids,
		Players: players,
	}, nil
}

// calculateMoonValue 计算月相值
func calculateMoonValue(card *card.Card) float64 {
	// TODO: 实现月相值计算规则
	return 1.0
}
