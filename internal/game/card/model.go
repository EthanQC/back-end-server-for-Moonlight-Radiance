package card

// BasicCard 基础牌
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

// SkillCard 功能牌
type SkillCard string

const ()

// CardType 卡牌类型
type CardType int

const (
	BasicCardType CardType = 1 //基础牌
	SkillCardType CardType = 2 //功能牌
)

// Card 卡牌基础结构
type Card struct {
	ID          uint     `gorm:"primarykey"`
	Name        string   `gorm:"size:50;not null"`
	Type        CardType `gorm:"not null"`
	Cost        int      `gorm:"not null"`
	Description string   `gorm:"size:500;not null"`
}

// PlayerCardState 玩家在对局中的卡牌状态
type PlayerCardState struct {
	ID              uint   `gorm:"primarykey"`
	GameID          uint   `gorm:"not null"`               // 对局ID
	PlayerID        uint   `gorm:"not null"`               // 玩家ID
	HandCardIDs     []uint `gorm:"type:json"`              // 手牌ID列表
	DeckCardIDs     []uint `gorm:"type:json"`              // 牌库ID列表
	DiscardCardIDs  []uint `gorm:"type:json"`              // 弃牌堆ID列表
	HandBasicCount  int    `gorm:"not null;default:0"`     // 手上的基础牌数量
	HandSkillCount  int    `gorm:"not null;default:0"`     // 手上的功能牌数量
	DeckBasicCount  int    `gorm:"not null;default:0"`     // 牌库的基础牌数量
	DeckSkillCount  int    `gorm:"not null;default:0"`     // 牌库的功能牌数量
	BasicCardPlayed bool   `gorm:"not null;default:false"` // 本回合是否出过基础牌
}

// CardStateResponse 返回给前端的卡牌状态
type CardStateResponse struct {
	Self struct {
		HandCards  []Card `json:"hand_cards"` // 手牌详细信息
		DeckCounts struct {
			Basic int `json:"basic"`
			Skill int `json:"skill"`
		} `json:"deck_counts"`
		DiscardCounts struct {
			Basic int `json:"basic"`
			Skill int `json:"skill"`
		} `json:"discard_counts"`
	} `json:"self"`
	Opponent struct {
		HandCounts struct {
			Basic int `json:"basic"`
			Skill int `json:"skill"`
		} `json:"hand_counts"`
		DeckCounts struct {
			Basic int `json:"basic"`
			Skill int `json:"skill"`
		} `json:"deck_counts"`
		DiscardCounts struct {
			Basic int `json:"basic"`
			Skill int `json:"skill"`
		} `json:"discard_counts"`
	} `json:"opponent"`
}
