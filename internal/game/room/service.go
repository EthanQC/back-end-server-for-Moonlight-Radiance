package room

import (
	"errors"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/internal/game/battlemap"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/internal/game/card"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/internal/game/racemap"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/common"
	"gorm.io/gorm"
)

type RoomService struct {
	db               *gorm.DB
	battlemapService *battlemap.BattleMapService
	racemapService   *racemap.RaceMapService
	cardService      *card.CardService
}

func NewRoomService() *RoomService {
	return &RoomService{
		db:               common.DB,
		battlemapService: battlemap.NewBattleMapService(),
		racemapService:   racemap.NewRaceMapService(),
		cardService:      card.NewCardService(),
	}
}

// CreateRoom 创建房间
func (s *RoomService) CreateRoom(playerID uint) (*Room, error) {
	var room *Room
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 创建对战地图
		battleMap, err := s.battlemapService.CreateMap(playerID)
		if err != nil {
			return err
		}

		// 2. 创建竞速地图
		raceMap, err := s.racemapService.CreateMap()
		if err != nil {
			return err
		}

		// 3. 创建房间
		newRoom := &Room{
			Status:          Waiting,
			Player1ID:       playerID,
			CurrentPlayerID: playerID,
			BattleMapID:     battleMap.ID,
			RaceMapID:       raceMap.ID,
		}

		if err := tx.Create(newRoom).Error; err != nil {
			return err
		}

		// 4. 初始化玩家进度
		progress := &PlayerProgress{
			RoomID:    newRoom.ID,
			PlayerID:  playerID,
			Position:  0,
			MoonValue: 0,
		}

		if err := tx.Create(progress).Error; err != nil {
			return err
		}

		// 5. 初始化玩家牌组
		if err := s.cardService.InitializePlayerDeck(newRoom.ID, playerID); err != nil {
			return err
		}

		room = newRoom
		return nil
	})

	if err != nil {
		return nil, err
	}

	return room, nil
}

// JoinRoom 加入房间
func (s *RoomService) JoinRoom(roomID uint, playerID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var room Room
		if err := tx.First(&room, roomID).Error; err != nil {
			return err
		}

		if room.Status != Waiting {
			return errors.New("room is not available")
		}

		if room.Player1ID == playerID {
			return errors.New("you are already in this room")
		}

		// 更新房间状态
		room.Player2ID = &playerID
		room.Status = Playing
		if err := tx.Save(&room).Error; err != nil {
			return err
		}

		// 创建玩家进度
		progress := &PlayerProgress{
			RoomID:    roomID,
			PlayerID:  playerID,
			Position:  0,
			MoonValue: 0,
		}

		if err := tx.Create(progress).Error; err != nil {
			return err
		}

		// 初始化玩家牌组
		return s.cardService.InitializePlayerDeck(roomID, playerID)
	})
}

// EndTurn 结束回合
func (s *RoomService) EndTurn(roomID uint, playerID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var room Room
		if err := tx.First(&room, roomID).Error; err != nil {
			return err
		}

		if room.CurrentPlayerID != playerID {
			return errors.New("not your turn")
		}

		// 1. 结束当前玩家的回合
		if err := s.cardService.EndTurn(roomID, playerID); err != nil {
			return err
		}

		// 2. 切换当前玩家
		if *room.Player2ID == playerID {
			room.CurrentPlayerID = room.Player1ID
		} else {
			room.CurrentPlayerID = *room.Player2ID
		}

		return tx.Save(&room).Error
	})
}

// GetRoomState 获取房间状态
func (s *RoomService) GetRoomState(roomID uint) (*Room, error) {
	var room Room
	if err := s.db.First(&room, roomID).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

// GetPlayerProgress 获取玩家进度
func (s *RoomService) GetPlayerProgress(roomID uint, playerID uint) (*PlayerProgress, error) {
	var progress PlayerProgress
	if err := s.db.Where("room_id = ? AND player_id = ?", roomID, playerID).First(&progress).Error; err != nil {
		return nil, err
	}
	return &progress, nil
}
