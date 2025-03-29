package card

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CardHandler struct {
	service *CardService
}

func NewCardHandler(s *CardService) *CardHandler {
	return &CardHandler{service: s}
}

// InitDeckHandler 初始化牌组
func (h *CardHandler) InitDeckHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		GameID uint `json:"game_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := h.service.InitializePlayerDeck(req.GameID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "牌组初始化成功"})
}

// GetCardStateHandler 获取卡牌状态
func (h *CardHandler) GetCardStateHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		GameID uint `json:"game_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	state, err := h.service.GetCardState(req.GameID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, state)
}

// DrawCardsHandler 抽初始手牌
func (h *CardHandler) DrawCardsHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		GameID uint `json:"game_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := h.service.DrawInitialCards(req.GameID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取更新后的状态
	state, err := h.service.GetCardState(req.GameID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, state)
}

// PlayCardHandler 出牌
func (h *CardHandler) PlayCardHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		GameID uint `json:"game_id" binding:"required"`
		CardID uint `json:"card_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := h.service.PlayCard(req.GameID, userID, req.CardID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取更新后的状态
	state, err := h.service.GetCardState(req.GameID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, state)
}

// EndTurnHandler 结束回合
func (h *CardHandler) EndTurnHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		GameID uint `json:"game_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := h.service.EndTurn(req.GameID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取更新后的状态
	state, err := h.service.GetCardState(req.GameID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, state)
}
