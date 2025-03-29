package racemap

import (
	"errors"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/common"
	"gorm.io/gorm"
)

const (
	DefaultMapLength = 20 // 默认地图长度
)

type RaceMapService struct {
	db *gorm.DB
}

func NewRaceMapService() *RaceMapService {
	return &RaceMapService{
		db: common.DB,
	}
}

// CreateMap 创建新的竞速地图
func (s *RaceMapService) CreateMap() (*RaceMap, error) {
	raceMap := &RaceMap{
		Length: DefaultMapLength,
	}

	if err := s.db.Create(raceMap).Error; err != nil {
		return nil, err
	}

	return raceMap, nil
}

// InitPlayerPosition 初始化玩家位置
func (s *RaceMapService) InitPlayerPosition(mapID uint, playerID uint) error {
	position := &Position{
		RaceMapID: mapID,
		PlayerID:  playerID,
		Location:  0,
		MoonValue: 0,
	}

	return s.db.Create(position).Error
}

// MoveForward 向前移动
func (s *RaceMapService) MoveForward(mapID uint, playerID uint, moonValue float64) (*MoveResult, error) {
	var result *MoveResult

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 获取地图信息
		var raceMap RaceMap
		if err := tx.First(&raceMap, mapID).Error; err != nil {
			return err
		}

		// 获取当前位置
		var position Position
		if err := tx.Where("race_map_id = ? AND player_id = ?", mapID, playerID).First(&position).Error; err != nil {
			return err
		}

		// 计算新位置
		newLocation := position.Location + 1
		if newLocation > raceMap.Length {
			return errors.New("cannot move beyond map end")
		}

		// 更新位置和月光值
		position.Location = newLocation
		position.MoonValue += moonValue

		if err := tx.Save(&position).Error; err != nil {
			return err
		}

		result = &MoveResult{
			NewLocation: newLocation,
			MoonValue:   position.MoonValue,
			IsFinished:  newLocation == raceMap.Length,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetPosition 获取玩家位置
func (s *RaceMapService) GetPosition(mapID uint, playerID uint) (*Position, error) {
	var position Position
	if err := s.db.Where("race_map_id = ? AND player_id = ?", mapID, playerID).First(&position).Error; err != nil {
		return nil, err
	}
	return &position, nil
}
