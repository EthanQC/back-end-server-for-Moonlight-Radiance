package card

import (
	"encoding/json"
	"errors"
	"math/rand"
	"slices"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CardService struct {
	db *gorm.DB
}

func NewCardService(db *gorm.DB) *CardService {
	return &CardService{db: db}
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

// GetGamePlayers 获取游戏中的所有玩家
func (s *CardService) GetGamePlayers(gameID uint) ([]uint, error) {
	var players []struct {
		PlayerID uint
	}
	if err := s.db.Table("game_players").
		Where("game_id = ?", gameID).
		Order("position").
		Find(&players).Error; err != nil {
		return nil, err
	}

	playerIDs := make([]uint, len(players))
	for i, p := range players {
		playerIDs[i] = p.PlayerID
	}
	return playerIDs, nil
}

// GetCardState 获取玩家的卡牌状态
func (s *CardService) GetCardState(gameID, playerID uint) (*CardStateResponse, error) {
	// 获取玩家的状态
	var state PlayerCardState
	if err := s.db.Where("game_id = ? AND player_id = ?", gameID, playerID).First(&state).Error; err != nil {
		return nil, err
	}

	// 获取所有玩家ID
	playerIDs, err := s.GetGamePlayers(gameID)
	if err != nil {
		return nil, err
	}

	// 构建响应
	response := &CardStateResponse{}

	// 处理自己的手牌
	if err := s.processOwnCards(&state, response); err != nil {
		return nil, err
	}

	// 处理所有对手的状态
	response.Opponents = make([]OpponentState, 0)
	for pos, pid := range playerIDs {
		if pid == playerID {
			continue
		}

		var oppState PlayerCardState
		if err := s.db.Where("game_id = ? AND player_id = ?", gameID, pid).
			First(&oppState).Error; err != nil {
			return nil, err
		}

		opponent := OpponentState{
			PlayerID: pid,
			Position: pos + 1,
		}
		opponent.HandCounts.Basic = oppState.HandBasicCount
		opponent.HandCounts.Skill = oppState.HandSkillCount
		opponent.DeckCounts.Basic = oppState.DeckBasicCount
		opponent.DeckCounts.Skill = oppState.DeckSkillCount

		// 统计对手弃牌堆
		if err := s.processOpponentDiscards(&oppState, &opponent); err != nil {
			return nil, err
		}

		response.Opponents = append(response.Opponents, opponent)
	}

	response.Stage = state.Stage
	return response, nil
}

// processOwnCards 处理玩家自己的卡牌
func (s *CardService) processOwnCards(state *PlayerCardState, response *CardStateResponse) error {
	// 处理手牌
	var handCardIDs []uint
	if err := json.Unmarshal(state.HandCardIDs, &handCardIDs); err != nil {
		return err
	}
	if len(handCardIDs) > 0 {
		if err := s.db.Find(&response.Self.HandCards, handCardIDs).Error; err != nil {
			return err
		}
	}

	// 设置牌库数量
	response.Self.DeckCounts.Basic = state.DeckBasicCount
	response.Self.DeckCounts.Skill = state.DeckSkillCount

	// 统计弃牌堆
	return s.processDiscardPile(state.DiscardCardIDs, &response.Self.DiscardCounts)
}

// processOpponentDiscards 处理对手的弃牌堆
func (s *CardService) processOpponentDiscards(state *PlayerCardState, opponent *OpponentState) error {
	return s.processDiscardPile(state.DiscardCardIDs, &opponent.DiscardCounts)
}

// processDiscardPile 处理弃牌堆统计
func (s *CardService) processDiscardPile(discardIDs datatypes.JSON, counts interface{}) error {
	var cardIDs []uint
	if err := json.Unmarshal(discardIDs, &cardIDs); err != nil {
		return err
	}

	if len(cardIDs) > 0 {
		var cards []Card
		if err := s.db.Find(&cards, cardIDs).Error; err != nil {
			return err
		}

		// 使用类型断言处理不同的计数结构
		switch c := counts.(type) {
		case *struct{ Basic, Skill int }:
			for _, card := range cards {
				if card.Type == BasicCardType {
					c.Basic++
				} else {
					c.Skill++
				}
			}
		}
	}
	return nil
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
