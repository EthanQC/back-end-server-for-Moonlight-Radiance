package battlemap

// Grid 地图格子
type Grid struct {
	Index     int     `json:"index"`      // 格子索引
	CardID    *uint   `json:"card_id"`    // 放置的卡牌ID
	OwnerID   *uint   `json:"owner_id"`   // 归属玩家ID
	MoonValue float64 `json:"moon_value"` // 当前月光值
}

// GameMap 游戏地图
type GameMap struct {
	ID     uint   `gorm:"primarykey"`
	GameID uint   `gorm:"not null"`        // 对局ID
	Size   int    `gorm:"not null"`        // 地图大小
	Grids  []Grid `gorm:"serializer:json"` // 格子状态
}

// MapStateResponse 返回给前端的地图状态
type MapStateResponse struct {
	MapID     uint   `json:"map_id"`
	Size      int    `json:"size"`
	Grids     []Grid `json:"grids"`
	MoonValue struct {
		Self     float64 `json:"self"`     // 自己的月光值
		Opponent float64 `json:"opponent"` // 对手的月光值
	} `json:"moon_value"`
}
