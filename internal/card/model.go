package card

import "gorm.io/datatypes"

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

const (
	MoonProphecy   SkillCard = "月之预言"
	MarsfallImpact SkillCard = "荧惑坠月"
	CassiaAegis    SkillCard = "桂魄含光"
)

// CardType 卡牌类型
type CardType int

const (
	BasicCardType CardType = 1 //基础牌
	SkillCardType CardType = 2 //功能牌
)

const (
	InitialBasicCards    = 2 // 初始抽取基础牌数量
	InitialSkillCards    = 1 // 初始抽取功能牌数量
	MaxBasicCards        = 3 // 基础牌手牌上限
	BaseSkillCards       = 3 // 功能牌基础上限
	BasicDeckSize        = 8 // 牌库基础牌数量
	InitialSkillDeckSize = 4 // 初始牌库功能牌数量
)

// PlayerState 玩家的游戏进程状态
type PlayerState int

const (
	Stage1 PlayerState = iota // 初始阶段，功能牌上限3
	Stage2                    // 功能牌上限3
	Stage3                    // 功能牌上限4
	Stage4                    // 功能牌上限5
	Stage5                    // 功能牌上限6
	Stage6                    // 功能牌上限7
)

// GetMaxSkillCards 获取当前阶段的功能牌上限
func (s PlayerState) GetMaxSkillCards() int {
	switch s {
	case Stage1, Stage2:
		return 3
	case Stage3:
		return 4
	case Stage4:
		return 5
	case Stage5:
		return 6
	case Stage6:
		return 7
	default:
		return BaseSkillCards
	}
}

// Card 卡牌基础结构
type Card struct {
	ID          uint
	Name        string
	Type        CardType
	Cost        int
	Description string
}

// PlayerCardState 玩家在对局中的卡牌状态
type PlayerCardState struct {
	ID              uint
	GameID          uint           // 对局ID
	PlayerID        uint           // 玩家ID
	Stage           PlayerState    // 游戏进程阶段
	HandCardIDs     datatypes.JSON // 手牌ID列表
	DeckCardIDs     datatypes.JSON // 牌库ID列表
	DiscardCardIDs  datatypes.JSON // 弃牌堆ID列表
	HandBasicCount  int            // 手上的基础牌数量
	HandSkillCount  int            // 手上的功能牌数量
	DeckBasicCount  int            // 牌库的基础牌数量
	DeckSkillCount  int            // 牌库的功能牌数量
	BasicCardPlayed bool           // 本回合是否出过基础牌
}

// TableName 指定表名
func (PlayerCardState) TableName() string {
	return "player_card_states"
}

// CardStateResponse 返回给前端的卡牌状态
type CardStateResponse struct {
	Self struct {
		HandCards  []Card `json:"handCards"` // 手牌详细信息
		DeckCounts struct {
			Basic int `json:"basic"`
			Skill int `json:"skill"`
		} `json:"deckCounts"`
		DiscardCounts struct {
			Basic int `json:"basic"`
			Skill int `json:"skill"`
		} `json:"discardCounts"`
	} `json:"self"`
	Opponent struct {
		HandCounts struct {
			Basic int `json:"basic"`
			Skill int `json:"skill"`
		} `json:"handCounts"`
		DeckCounts struct {
			Basic int `json:"basic"`
			Skill int `json:"skill"`
		} `json:"deckCounts"`
		DiscardCounts struct {
			Basic int `json:"basic"`
			Skill int `json:"skill"`
		} `json:"discardCounts"`
	} `json:"opponent"`
	Stage PlayerState `json:"stage"`
}
