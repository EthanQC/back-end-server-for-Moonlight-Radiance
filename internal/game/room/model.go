package room

import (
	"time"
)

// RoomStatus 房间状态
type RoomStatus int

const (
	Waiting  RoomStatus = iota // 等待玩家加入
	Playing                    // 对战中
	Finished                   // 对战结束
)

// Room 对战房间
type Room struct {
	ID              uint       `gorm:"primarykey"`
	Status          RoomStatus `gorm:"not null"`
	Player1ID       uint       `gorm:"not null"`     // 玩家1
	Player2ID       *uint      `gorm:"default:null"` // 玩家2(可能为空)
	CurrentPlayerID uint       `gorm:"not null"`     // 当前回合玩家
	BattleMapID     uint       `gorm:"not null"`     // 对战地图ID
	RaceMapID       uint       `gorm:"not null"`     // 竞速地图ID
	CreatedAt       time.Time  `gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime"`
}

// PlayerProgress 玩家在竞速地图上的进度
type PlayerProgress struct {
	RoomID    uint    `gorm:"primarykey;autoIncrement:false"`
	PlayerID  uint    `gorm:"primarykey;autoIncrement:false"`
	Position  int     `gorm:"not null"` // 当前位置
	MoonValue float64 `gorm:"not null"` // 累计月光值
}
