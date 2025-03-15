package card

// GridStatus 格子状态
type GridStatus struct {
	Index   int   `json:"index"`    // 格子索引
	Card    *Card `json:"card"`     // 放置的基础牌
	OwnerID uint  `json:"owner_id"` // 归属玩家ID
}

// GameMap 游戏地图
type GameMap struct {
	ID    uint         `gorm:"primarykey"`
	Type  string       `gorm:"size:20"` // PVE/PVP
	Size  int          `json:"size"`    // 地图大小
	Grids []GridStatus `gorm:"serializer:json"`
}

// GameState 游戏状态
type GameState struct {
	ID        uint    `gorm:"primarykey"`
	UserID    uint    `gorm:"not null"`
	MapID     uint    `gorm:"not null"`
	MoonValue float64 `gorm:"default:0"`       // 月光值
	BasicHand []Card  `gorm:"serializer:json"` // 基础牌手牌
	SkillHand []Card  `gorm:"serializer:json"` // 功能牌手牌
}
