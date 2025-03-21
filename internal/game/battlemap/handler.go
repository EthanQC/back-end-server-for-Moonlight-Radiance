package battlemap

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BattleMapHandler struct {
	service *BattleMapService
}

func NewBattleMapHandler() *BattleMapHandler {
	return &BattleMapHandler{
		service: NewBattleMapService(),
	}
}

// CreateMapHandler 创建对战地图
func (h *BattleMapHandler) CreateMapHandler(c *gin.Context) {
	var req struct {
		GameID uint `json:"game_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	battleMap, err := h.service.CreateMap(req.GameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, battleMap)
}

// PlaceCardHandler 放置卡牌
func (h *BattleMapHandler) PlaceCardHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		MapID  uint `json:"map_id" binding:"required"`
		X      int  `json:"x" binding:"required"`
		Y      int  `json:"y" binding:"required"`
		CardID uint `json:"card_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := h.service.PlaceCard(req.MapID, userID, req.X, req.Y, req.CardID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	state, err := h.service.GetMapState(req.MapID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, state)
}

// GetMapStateHandler 获取地图状态
func (h *BattleMapHandler) GetMapStateHandler(c *gin.Context) {
	var req struct {
		MapID uint `json:"map_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	state, err := h.service.GetMapState(req.MapID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, state)
}
