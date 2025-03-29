package room

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	service *RoomService
}

func NewRoomHandler(s *RoomService) *RoomHandler {
	return &RoomHandler{service: s}
}

// POST /api/room/create
// body: {"capacity": 2 or 3 or 4}
func (h *RoomHandler) CreateRoomHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		Capacity int `json:"capacity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	room, err := h.service.CreateRoom(userID, req.Capacity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"room_id":  room.ID,
		"host_id":  room.HostID,
		"capacity": room.Capacity,
		"status":   room.Status,
	})
}

// POST /api/room/join
// body: {"room_id": 123}
func (h *RoomHandler) JoinRoomHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		RoomID uint `json:"room_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.JoinRoom(req.RoomID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "joined room"})
}

// GET /api/room/state?room_id=xxx
func (h *RoomHandler) GetRoomStateHandler(c *gin.Context) {
	rid := c.Query("room_id")
	if rid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room_id required"})
		return
	}

	// 将房间 ID 转换为整数
	roomID, _ := strconv.Atoi(rid)
	state, err := h.service.GetRoomState(uint(roomID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, state)
}
