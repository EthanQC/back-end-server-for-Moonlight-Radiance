package card

import (
	"errors"
	"math/rand"
	"slices"
	"time"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/common"
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

	// 初始化牌库
	for _, card := range basicCards {
		state.DeckCardIDs = append(state.DeckCardIDs, card.ID)
		state.DeckBasicCount++
	}
	for _, card := range skillCards {
		state.DeckCardIDs = append(state.DeckCardIDs, card.ID)
		state.DeckSkillCount++
	}

	// 打乱牌库
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(state.DeckCardIDs), func(i, j int) {
		state.DeckCardIDs[i], state.DeckCardIDs[j] = state.DeckCardIDs[j], state.DeckCardIDs[i]
	})

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
	if len(state.HandCardIDs) > 0 {
		var handCards []Card
		if err := s.db.Where("id IN ?", state.HandCardIDs).Find(&handCards).Error; err != nil {
			return nil, err
		}
		response.Self.HandCards = handCards
	}

	// 统计弃牌堆
	if len(state.DiscardCardIDs) > 0 {
		var discardCards []Card
		if err := s.db.Where("id IN ?", state.DiscardCardIDs).Find(&discardCards).Error; err != nil {
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
		var opponentDiscardCards []Card
		if err := s.db.Where("id IN ?", opponentState.DiscardCardIDs).Find(&opponentDiscardCards).Error; err != nil {
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

	drawnCards := make([]uint, 0)
	drawnBasicCount := 0
	drawnSkillCount := 0

	// 随机抽取指定数量的牌
	for _, cardID := range state.DeckCardIDs {
		var card Card
		if err := s.db.First(&card, cardID).Error; err != nil {
			return err
		}

		if card.Type == BasicCardType && drawnBasicCount < basicCount {
			drawnCards = append(drawnCards, cardID)
			drawnBasicCount++
		} else if card.Type == SkillCardType && drawnSkillCount < skillCount {
			drawnCards = append(drawnCards, cardID)
			drawnSkillCount++
		}

		if drawnBasicCount >= basicCount && drawnSkillCount >= skillCount {
			break
		}
	}

	// 更新状态
	state.HandCardIDs = append(state.HandCardIDs, drawnCards...)
	state.HandBasicCount += drawnBasicCount
	state.HandSkillCount += drawnSkillCount

	// 更新牌库
	newDeckCards := make([]uint, 0)
	for _, cardID := range state.DeckCardIDs {
		if !slices.Contains(drawnCards, cardID) {
			newDeckCards = append(newDeckCards, cardID)
		}
	}
	state.DeckCardIDs = newDeckCards
	state.DeckBasicCount -= drawnBasicCount
	state.DeckSkillCount -= drawnSkillCount

	return s.db.Save(&state).Error
}

// DrawInitialCards 初始抽牌
func (s *CardService) DrawInitialCards(gameID, playerID uint) error {
	return s.DrawCards(gameID, playerID, Basic, Skill)
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

	if !slices.Contains(state.HandCardIDs, cardID) {
		return errors.New("卡牌不在手牌中")
	}

	// 验证出牌规则
	if playedCard.Type == BasicCardType && state.BasicCardPlayed {
		return errors.New("本回合已经出过基础牌")
	}
	state.BasicCardPlayed = true

	// 从手牌中移除
	newHandCards := make([]uint, 0)
	for _, handCardID := range state.HandCardIDs {
		if handCardID != cardID {
			newHandCards = append(newHandCards, handCardID)
		}
	}
	state.HandCardIDs = newHandCards

	// 更新手牌数量
	if playedCard.Type == BasicCardType {
		state.HandBasicCount--
		state.BasicCardPlayed = true
	} else {
		state.HandSkillCount--
	}

	// 添加到弃牌堆
	state.DiscardCardIDs = append(state.DiscardCardIDs, cardID)

	return s.db.Save(&state).Error
}

// EndTurn 结束回合
func (s *CardService) EndTurn(gameID, playerID uint) error {
	var state PlayerCardState
	if err := s.db.Where("game_id = ? AND player_id = ?", gameID, playerID).First(&state).Error; err != nil {
		return err
	}

	// 重置回合状态
	state.BasicCardPlayed = false

	// 补充手牌
	basicNeed := 3 - state.HandBasicCount
	skillNeed := 3 - state.HandSkillCount

	if basicNeed > 0 || skillNeed > 0 {
		if err := s.DrawCards(gameID, playerID, basicNeed, skillNeed); err != nil {
			return err
		}
	}

	return s.db.Save(&state).Error
}
