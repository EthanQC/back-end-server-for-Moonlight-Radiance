package racemap

import (
	"time"
)

// RaceMap 竞速地图
type RaceMap struct {
	ID        uint      `gorm:"primarykey"`
	Length    int       `gorm:"not null"` // 地图总长度(格子数)
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// Position 玩家在竞速地图上的位置
type Position struct {
	RaceMapID uint    `gorm:"primarykey"`
	PlayerID  uint    `gorm:"primarykey"`
	Location  int     `gorm:"not null"` // 当前位置
	MoonValue float64 `gorm:"not null"` // 月光值
}

// MoveResult 移动结果
type MoveResult struct {
	NewLocation int     `json:"new_location"`
	MoonValue   float64 `json:"moon_value"`
	IsFinished  bool    `json:"is_finished"`
}
