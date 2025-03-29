package battlemap

import (
	"gorm.io/datatypes"
)

// CardPlacement 卡牌放置信息
type CardPlacement struct {
	CardID    uint    `json:"card_id"`    // 放置的卡牌ID
	PlayerID  uint    `json:"player_id"`  // 放置的玩家ID
	MoonValue float64 `json:"moon_value"` // 月相值
}

// Grid 地图格子
type Grid struct {
	Index     int            `json:"index"`
	Position  Position       `json:"position"`
	Placement *CardPlacement `json:"placement,omitempty"`
}

// Position 格子位置
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// BattleMap 对战地图
type BattleMap struct {
	ID        uint           `gorm:"primarykey"`
	GameID    uint           `gorm:"not null"`
	Size      int            `gorm:"not null"` // 标准大小8x8
	Grids     datatypes.JSON `gorm:"type:json"`
	CreatedAt int64          `gorm:"autoCreateTime"`
	UpdatedAt int64          `gorm:"autoUpdateTime"`
}

// MapState 地图状态
type MapState struct {
	MapID   uint          `json:"map_id"`
	Grids   []Grid        `json:"grids"`
	Players []PlayerState `json:"players"`
}

// PlayerState 玩家状态
type PlayerState struct {
	ID        uint    `json:"id"`
	MoonValue float64 `json:"moon_value"`
}
