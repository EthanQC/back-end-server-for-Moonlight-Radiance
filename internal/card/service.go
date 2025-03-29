package card

import (
	"encoding/json"
	"errors"
	"math/rand"
	"slices"
	"time"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/common"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CardService struct {
	db *gorm.DB
}

func NewCardService() *CardService {
	return &CardService{
		db: common.DB,
	}
}

// InitializePlayerDeck 初始化玩家牌组
func (s *CardService) InitializePlayerDeck(gameID, playerID uint) error {
	// 创建初始牌组状态
	state := PlayerCardState{
		GameID:   gameID,
		PlayerID: playerID,
	}

	// 获取所有基础牌
	var basicCards []Card
	if err := s.db.Where("type = ?", BasicCardType).Find(&basicCards).Error; err != nil {
		return err
	}

	// 获取所有功能牌
	var skillCards []Card
	if err := s.db.Where("type = ?", SkillCardType).Find(&skillCards).Error; err != nil {
		return err
	}

	// 创建牌库数组
	deckCards := make([]uint, 0)

	// 初始化牌库
	for _, card := range basicCards {
		deckCards = append(deckCards, card.ID)
		state.DeckBasicCount++
	}
	for _, card := range skillCards {
		deckCards = append(deckCards, card.ID)
		state.DeckSkillCount++
	}

	// 打乱牌库
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(state.DeckCardIDs), func(i, j int) {
		deckCards[i], deckCards[j] = deckCards[j], deckCards[i]
	})

	// 将数组转换为 JSON
	deckJSON, err := json.Marshal(deckCards)
	if err != nil {
		return err
	}
	state.DeckCardIDs = datatypes.JSON(deckJSON)

	// 初始化空的手牌和弃牌堆
	emptyJSON := datatypes.JSON("[]")
	state.HandCardIDs = emptyJSON
	state.DiscardCardIDs = emptyJSON

	return s.db.Create(&state).Error
}

// GetCardState 获取玩家的卡牌状态
func (s *CardService) GetCardState(gameID, playerID uint) (*CardStateResponse, error) {
	// 获取玩家的状态
	var state PlayerCardState
	if err := s.db.Where("game_id = ? AND player_id = ?", gameID, playerID).First(&state).Error; err != nil {
		return nil, err
	}

	// 获取对手的状态
	var opponentState PlayerCardState
	if err := s.db.Where("game_id = ? AND player_id != ?", gameID, playerID).First(&opponentState).Error; err != nil {
		return nil, err
	}

	// 构建响应
	response := &CardStateResponse{}

	// 获取并设置手牌详细信息
	var handCardIDs []uint
	if err := json.Unmarshal(state.HandCardIDs, &handCardIDs); err != nil {
		return nil, err
	}
	if len(handCardIDs) > 0 {
		if err := s.db.Find(&response.Self.HandCards, handCardIDs).Error; err != nil {
			return nil, err
		}
	}

	// 统计弃牌堆
	if len(state.DiscardCardIDs) > 0 {
		var discardCardIDs []uint
		if err := json.Unmarshal(state.DiscardCardIDs, &discardCardIDs); err != nil {
			return nil, err
		}

		var discardCards []Card
		if err := s.db.Find(&discardCards, discardCardIDs).Error; err != nil {
			return nil, err
		}

		for _, card := range discardCards {
			if card.Type == BasicCardType {
				response.Self.DiscardCounts.Basic++
			} else {
				response.Self.DiscardCounts.Skill++
			}
		}
	}

	// 设置牌库数量
	response.Self.DeckCounts.Basic = state.DeckBasicCount
	response.Self.DeckCounts.Skill = state.DeckSkillCount

	// 设置对手信息
	response.Opponent.HandCounts.Basic = opponentState.HandBasicCount
	response.Opponent.HandCounts.Skill = opponentState.HandSkillCount
	response.Opponent.DeckCounts.Basic = opponentState.DeckBasicCount
	response.Opponent.DeckCounts.Skill = opponentState.DeckSkillCount

	// 统计对手弃牌堆
	if len(opponentState.DiscardCardIDs) > 0 {
		var opponentDiscardIDs []uint
		if err := json.Unmarshal(opponentState.DiscardCardIDs, &opponentDiscardIDs); err != nil {
			return nil, err
		}

		var opponentDiscardCards []Card
		if err := s.db.Find(&opponentDiscardCards, opponentDiscardIDs).Error; err != nil {
			return nil, err
		}

		for _, card := range opponentDiscardCards {
			if card.Type == BasicCardType {
				response.Opponent.DiscardCounts.Basic++
			} else {
				response.Opponent.DiscardCounts.Skill++
			}
		}
	}

	return response, nil
}

// DrawCards 抽指定数量的牌
func (s *CardService) DrawCards(gameID, playerID uint, basicCount, skillCount int) error {
	var state PlayerCardState
	if err := s.db.Where("game_id = ? AND player_id = ?", gameID, playerID).First(&state).Error; err != nil {
		return err
	}

	// 获取牌库
	var deckCards []uint
	if err := json.Unmarshal(state.DeckCardIDs, &deckCards); err != nil {
		return err
	}

	// 打乱牌库
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(deckCards), func(i, j int) {
		deckCards[i], deckCards[j] = deckCards[j], deckCards[i]
	})

	// 明确需要抽取的牌数
	maxBasic := MaxBasicCards - state.HandBasicCount
	if basicCount > maxBasic {
		basicCount = maxBasic
	}
	maxSkill := state.Stage.GetMaxSkillCards() - state.HandSkillCount
	if skillCount > maxSkill {
		skillCount = maxSkill
	}

	// 抽牌
	drawnCards := make([]uint, 0)
	drawnBasicCount := 0
	drawnSkillCount := 0
	remainingDeck := make([]uint, 0)

	for _, cardID := range deckCards {
		var card Card
		if err := s.db.First(&card, cardID).Error; err != nil {
			return err
		}

		drawn := false
		if card.Type == BasicCardType && drawnBasicCount < basicCount {
			drawnCards = append(drawnCards, cardID)
			drawnBasicCount++
			drawn = true
		} else if card.Type == SkillCardType && drawnSkillCount < skillCount {
			drawnCards = append(drawnCards, cardID)
			drawnSkillCount++
			drawn = true
		}

		if !drawn {
			remainingDeck = append(remainingDeck, cardID)
		}

		if drawnBasicCount >= basicCount && drawnSkillCount >= skillCount {
			// 将剩余的牌加入牌库
			remainingDeck = append(remainingDeck, deckCards[len(drawnCards):]...)
			break
		}
	}

	// 更新手牌
	var handCards []uint
	if err := json.Unmarshal(state.HandCardIDs, &handCards); err != nil {
		return err
	}
	handCards = append(handCards, drawnCards...)
	handJSON, err := json.Marshal(handCards)
	if err != nil {
		return err
	}
	state.HandCardIDs = datatypes.JSON(handJSON)

	// 更新牌库
	deckJSON, err := json.Marshal(remainingDeck)
	if err != nil {
		return err
	}
	state.DeckCardIDs = datatypes.JSON(deckJSON)

	// 更新计数
	state.HandBasicCount += drawnBasicCount
	state.HandSkillCount += drawnSkillCount
	state.DeckBasicCount -= drawnBasicCount
	state.DeckSkillCount -= drawnSkillCount

	return s.db.Save(&state).Error
}

// DrawInitialCards 初始抽牌
func (s *CardService) DrawInitialCards(gameID, playerID uint) error {
	return s.DrawCards(gameID, playerID, InitialBasicCards, InitialSkillCards)
}

// PlayCard 打出一张牌
func (s *CardService) PlayCard(gameID, playerID, cardID uint) error {
	var state PlayerCardState
	if err := s.db.Where("game_id = ? AND player_id = ?", gameID, playerID).First(&state).Error; err != nil {
		return err
	}

	// 验证是否有这张牌
	var playedCard Card
	if err := s.db.First(&playedCard, cardID).Error; err != nil {
		return err
	}

	var handCards []uint
	if err := json.Unmarshal(state.HandCardIDs, &handCards); err != nil {
		return err
	}
	if !slices.Contains(handCards, cardID) {
		return errors.New("卡牌不在手牌中")
	}

	// 验证出牌规则
	if playedCard.Type == BasicCardType && state.BasicCardPlayed {
		return errors.New("本回合已经出过基础牌")
	}
	state.BasicCardPlayed = true

	// 从手牌中移除
	newHandCards := make([]uint, 0)
	var currentHandCards []uint
	if err := json.Unmarshal(state.HandCardIDs, &currentHandCards); err != nil {
		return err
	}
	for _, handCardID := range currentHandCards {
		if handCardID != cardID {
			newHandCards = append(newHandCards, handCardID)
		}
	}
	handJSON, err := json.Marshal(newHandCards)
	if err != nil {
		return err
	}
	state.HandCardIDs = datatypes.JSON(handJSON)

	// 更新手牌数量
	if playedCard.Type == BasicCardType {
		state.HandBasicCount--
		state.BasicCardPlayed = true
	} else {
		state.HandSkillCount--
	}

	// 添加到弃牌堆
	var discardCards []uint
	if err := json.Unmarshal(state.DiscardCardIDs, &discardCards); err != nil {
		return err
	}
	discardCards = append(discardCards, cardID)
	discardJSON, err := json.Marshal(discardCards)
	if err != nil {
		return err
	}
	state.DiscardCardIDs = datatypes.JSON(discardJSON)

	return s.db.Save(&state).Error
}

// EndTurn 结束回合
func (s *CardService) EndTurn(gameID, playerID uint) error {
	// 使用事务确保状态一致性
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var state PlayerCardState
	if err := tx.Where("game_id = ? AND player_id = ?", gameID, playerID).First(&state).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 重置回合状态
	state.BasicCardPlayed = false

	// 更新状态
	if err := tx.Save(&state).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}

	// 计算需要补充的牌数
	basicNeed := MaxBasicCards - state.HandBasicCount
	skillNeed := state.Stage.GetMaxSkillCards() - state.HandSkillCount

	// 补充手牌（在事务外执行）
	if basicNeed > 0 || skillNeed > 0 {
		return s.DrawCards(gameID, playerID, basicNeed, skillNeed)
	}

	return nil
}
