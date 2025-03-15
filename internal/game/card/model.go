package card

// BasicCard 基础月相牌类型
type BasicCard string

const (
	NewMoon        BasicCard = "新月"
	WaxingCrescent BasicCard = "蛾眉月"
	FirstQuarter   BasicCard = "上弦月"
	WaxingGibbous  BasicCard = "盈凸月"
	FullMoon       BasicCard = "满月"
	WaningGibbous  BasicCard = "亏凸月"
	LastQuarter    BasicCard = "下弦月"
	WaningCrescent BasicCard = "残月"
)

// SkillCard 功能牌类型
type SkillCard string

const ()

type CardType int

const (
	BasicCardType CardType = 1 //基础牌
	SkillCardType CardType = 2 //功能牌
)

// Card 卡牌基础结构
type Card struct {
	ID          uint     `gorm:"primarykey"`
	Name        string   `gorm:"size:50;not null"`
	Type        CardType `gorm:"size:20;not null;"`
	Description string   `gorm:"size:500"`
	Chapter     int      `gorm:"default:1"` // 解锁所需章节
}

// PlayerDeck 玩家牌组
type PlayerDeck struct {
	ID         uint   `gorm:"primarykey"`
	UserID     uint   `gorm:"not null"`
	Name       string `gorm:"size:50"`
	BasicCards int    `gorm:"default:0"` // 基础牌数量
	SkillCards int    `gorm:"default:0"` // 功能牌数量
	CardIDs    []uint `gorm:"type:json"` // 使用JSON类型存储卡牌ID数组
}
