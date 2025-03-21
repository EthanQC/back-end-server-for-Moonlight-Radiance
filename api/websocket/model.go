package websocket

// EventType 事件类型
type EventType string

const (
	// 游戏事件
	EventCardPlayed   EventType = "card_played"   // 出牌
	EventCardPlaced   EventType = "card_placed"   // 放置卡牌
	EventTurnChanged  EventType = "turn_changed"  // 回合切换
	EventGameStarted  EventType = "game_started"  // 游戏开始
	EventGameEnded    EventType = "game_ended"    // 游戏结束
	EventStateChanged EventType = "state_changed" // 状态更新

	// 地图事件
	EventGridOccupied EventType = "grid_occupied" // 格子被占领
	EventPlayerMoved  EventType = "player_moved"  // 玩家移动

	// 系统事件
	EventError       EventType = "error"       // 错误消息
	EventReconnected EventType = "reconnected" // 重连成功
)

// Event WebSocket事件
type Event struct {
	Type     EventType   `json:"type"`
	RoomID   uint        `json:"room_id"`
	PlayerID uint        `json:"player_id,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}
