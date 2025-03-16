package card

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardFlow(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	s := NewCardService()
	game_id := uint(1)
	player_id := uint(1)
	opponent_id := uint(2)

	// 1. 测试初始化牌组
	t.Run("Initialize Deck", func(t *testing.T) {
		// 先初始化对手的牌组
		err := s.InitializePlayerDeck(game_id, opponent_id)
		assert.NoError(t, err)

		// 再初始化玩家自己的牌组
		err = s.InitializePlayerDeck(game_id, player_id)
		assert.NoError(t, err)

		state, err := s.GetCardState(game_id, player_id)
		assert.NoError(t, err)
		assert.Equal(t, BasicDeckSize, state.Self.DeckCounts.Basic)
		assert.Equal(t, InitialSkillDeckSize, state.Self.DeckCounts.Skill)
	})

	// 2. 测试初始抽牌
	t.Run("Draw Initial Cards", func(t *testing.T) {
		err := s.DrawInitialCards(game_id, player_id)
		assert.NoError(t, err)

		state, err := s.GetCardState(game_id, player_id)
		assert.NoError(t, err)
		assert.Len(t, state.Self.HandCards, InitialBasicCards+InitialSkillCards) // 验证手牌数量

		// 检查手牌类型
		basicCount := 0
		skillCount := 0
		for _, card := range state.Self.HandCards {
			if card.Type == BasicCardType {
				basicCount++
			} else {
				skillCount++
			}
		}
		assert.Equal(t, InitialBasicCards, basicCount, "应该有2张基础牌")
		assert.Equal(t, InitialSkillCards, skillCount, "应该有1张功能牌")
	})

	// 3. 测试出牌
	t.Run("Play Card", func(t *testing.T) {
		state, err := s.GetCardState(game_id, player_id)
		assert.NoError(t, err)

		// 出一张基础牌
		if len(state.Self.HandCards) > 0 {
			// 出一张基础牌
			var basicCard Card
			for _, card := range state.Self.HandCards {
				if card.Type == BasicCardType {
					basicCard = card
					break
				}
			}

			initialHandCount := len(state.Self.HandCards)
			err = s.PlayCard(game_id, player_id, basicCard.ID)
			assert.NoError(t, err)

			newState, err := s.GetCardState(game_id, player_id)
			assert.NoError(t, err)
			assert.Equal(t, initialHandCount-1, len(newState.Self.HandCards))

			// 尝试在同一回合再出一张基础牌
			for _, card := range newState.Self.HandCards {
				if card.Type == BasicCardType {
					err = s.PlayCard(game_id, player_id, card.ID)
					assert.Error(t, err)
					break
				}
			}
		}
	})

	// 4. 测试结束回合
	t.Run("End Turn", func(t *testing.T) {
		// 获取结束回合前的状态
		beforeState, err := s.GetCardState(game_id, player_id)
		assert.NoError(t, err)

		err = s.EndTurn(game_id, player_id)
		assert.NoError(t, err)

		afterState, err := s.GetCardState(game_id, player_id)
		assert.NoError(t, err)

		// 统计手牌类型
		basicCount := 0
		skillCount := 0
		for _, card := range afterState.Self.HandCards {
			if card.Type == BasicCardType {
				basicCount++
			} else {
				skillCount++
			}
		}

		// 验证基础牌和功能牌的数量
		assert.Equal(t, MaxBasicCards, basicCount, "结束回合后应该有3张基础牌")
		maxSkill := beforeState.Stage.GetMaxSkillCards()
		assert.Equal(t, maxSkill, skillCount, "结束回合后应该有3张功能牌")

		// 验证总手牌数
		assert.Equal(t, MaxBasicCards+maxSkill, len(afterState.Self.HandCards),
			"结束回合后总手牌数应该为基础牌上限+功能牌上限")
	})

	// 5. 测试对手视角
	t.Run("Opponent View", func(t *testing.T) {
		// 初始化对手牌组
		err := s.InitializePlayerDeck(game_id, opponent_id)
		assert.NoError(t, err)

		// 获取对手视角的状态
		state, err := s.GetCardState(game_id, opponent_id)
		assert.NoError(t, err)
		assert.NotNil(t, state.Opponent.HandCounts)
		assert.Equal(t, 0, len(state.Self.HandCards))
	})
}
