package racemap

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RaceMapHandler struct {
	service *RaceMapService
}

func NewRaceMapHandler() *RaceMapHandler {
	return &RaceMapHandler{
		service: NewRaceMapService(),
	}
}

// CreateMapHandler 创建竞速地图
func (h *RaceMapHandler) CreateMapHandler(c *gin.Context) {
	raceMap, err := h.service.CreateMap()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, raceMap)
}

// MoveForwardHandler 向前移动
func (h *RaceMapHandler) MoveForwardHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		MapID     uint    `json:"map_id" binding:"required"`
		MoonValue float64 `json:"moon_value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	result, err := h.service.MoveForward(req.MapID, userID, req.MoonValue)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetPositionHandler 获取位置信息
func (h *RaceMapHandler) GetPositionHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		MapID uint `json:"map_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	position, err := h.service.GetPosition(req.MapID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, position)
}
