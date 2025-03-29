package room

import (
	"errors"

	"gorm.io/gorm"
)

type RoomService struct {
	db *gorm.DB
}

func NewRoomService(db *gorm.DB) *RoomService {
	return &RoomService{db: db}
}

// CreateRoom 创建房间 (带房主)
func (s *RoomService) CreateRoom(hostID uint, capacity int) (*Room, error) {
	if capacity < 2 || capacity > 4 {
		return nil, errors.New("capacity must be between 2 and 4")
	}

	newRoom := &Room{
		HostID:   hostID,
		Capacity: capacity,
		Status:   Waiting,
	}

	// 启用事务，接收一个匿名函数作为函数，函数的参数是一个 *gorm.DB 类型的事务对象
	// 该函数会在事务开始时被调用，并且在事务结束时被提交或回滚
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 插入 Room
		if err := tx.Create(newRoom).Error; err != nil {
			return err
		}

		// 2. 插入房主到 room_players
		rp := &RoomPlayer{
			RoomID:   newRoom.ID,
			PlayerID: hostID,
		}
		if err := tx.Create(rp).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return newRoom, nil
}

// JoinRoom 加入现有房间
func (s *RoomService) JoinRoom(roomID, userID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var room Room
		if err := tx.First(&room, roomID).Error; err != nil {
			return err
		}

		if room.Status != Waiting {
			return errors.New("room not waiting")
		}

		// 检查房间人数
		var count int64
		if err := tx.Model(&RoomPlayer{}).Where("room_id = ?", roomID).Count(&count).Error; err != nil {
			return err
		}
		if int(count) >= room.Capacity {
			return errors.New("room is full")
		}

		// 检查是否已经在房间
		var cnt int64
		if err := tx.Model(&RoomPlayer{}).Where("room_id = ? AND user_id = ?", roomID, userID).Count(&cnt).Error; err != nil {
			return err
		}
		if cnt > 0 {
			return errors.New("already in room")
		}

		// 加入
		rp := RoomPlayer{
			RoomID:   roomID,
			PlayerID: userID,
		}
		if err := tx.Create(&rp).Error; err != nil {
			return err
		}

		// 如果人数满 -> 状态改为Playing
		count++
		if int(count) == room.Capacity {
			room.Status = Playing
			if err := tx.Save(&room).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *RoomService) GetRoomState(roomID uint) (*RoomStateResponse, error) {
	var room Room
	if err := s.db.First(&room, roomID).Error; err != nil {
		return nil, err
	}

	var response []RoomPlayer
	if err := s.db.Where("room_id = ?", roomID).Find(&response).Error; err != nil {
		return nil, err
	}
	playerIDs := make([]uint, 0, len(response))
	for _, rp := range response {
		playerIDs = append(playerIDs, rp.PlayerID)
	}
	return &RoomStateResponse{
		RoomID:   room.ID,
		HostID:   room.HostID,
		PlayerID: playerIDs,
		Capacity: room.Capacity,
		Status:   room.Status,
	}, nil
}
