package game

import (
	"time"
)

type GameStaus int

const (
	Runing GameStaus = iota
	Ended
)

type Game struct {
	ID             uint      `gorm:"primarykey"`
	RoomID         uint      `gorm:"not null"` // 房间ID
	Status         GameStaus `gorm:"not null"` // 游戏状态
	BattleMapID    uint      `gorm:"not null"` // 战斗地图ID
	RaceMapID      uint      `gorm:"not null"` // 竞速地图ID
	CurrntPlayerID uint      `gorm:"not null"` // 当前玩家ID
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}
