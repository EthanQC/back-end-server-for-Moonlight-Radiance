package room

import (
	"time"
)

// RoomStatus 房间状态
type RoomStatus int

const (
	Waiting  RoomStatus = iota // 等待玩家加入，0
	Playing                    // 对战中，1
	Finished                   // 对战结束，2
)

// Room 对战房间
type Room struct {
	ID        uint
	Status    RoomStatus
	Capacity  int  // 房间容量，2-4人
	HostID    uint // 房主ID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// RoomPlayer 房间玩家
type RoomPlayer struct {
	ID        uint
	RoomID    uint  // 房间ID
	Room      *Room // 关联房间
	PlayerID  uint  // 玩家ID
	IsReady   bool  // 是否准备
	CreatedAt time.Time
	UpdatedAt time.Time
}

// RoomStateResponse 返回房间信息 + 房间玩家
type RoomStateResponse struct {
	RoomID   uint       `json:"room_id"`
	HostID   uint       `json:"host_id"`
	PlayerID []uint     `json:"player_id"`
	Capacity int        `json:"capacity"`
	Status   RoomStatus `json:"status"`
}
