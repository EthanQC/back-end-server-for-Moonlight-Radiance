package room

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	service *RoomService
}

func NewRoomHandler() *RoomHandler {
	return &RoomHandler{
		service: NewRoomService(),
	}
}

// CreateRoomHandler 创建房间
func (h *RoomHandler) CreateRoomHandler(c *gin.Context) {
	userID := c.GetUint("user_id")

	room, err := h.service.CreateRoom(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, room)
}

// JoinRoomHandler 加入房间
func (h *RoomHandler) JoinRoomHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		RoomID uint `json:"room_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := h.service.JoinRoom(req.RoomID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "成功加入房间"})
}

// EndTurnHandler 结束回合
func (h *RoomHandler) EndTurnHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		RoomID uint `json:"room_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := h.service.EndTurn(req.RoomID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "回合结束"})
}

// GetRoomStateHandler 获取房间状态
func (h *RoomHandler) GetRoomStateHandler(c *gin.Context) {
	var req struct {
		RoomID uint `json:"room_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	state, err := h.service.GetRoomState(req.RoomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, state)
}
