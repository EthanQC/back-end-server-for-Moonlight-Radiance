package websocket

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境需要更严格的检查
	},
}

type Handler struct {
	hub *Hub
}

func NewHandler() *Handler {
	return &Handler{
		hub: NewHub(),
	}
}

// HandleConnection WebSocket连接处理
func (h *Handler) HandleConnection(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	playerID := c.GetUint("user_id")
	roomID, err := strconv.ParseUint(c.Query("room_id"), 10, 32)
	if err != nil {
		ws.Close()
		return
	}

	conn := newConnection(ws, h.hub, playerID, uint(roomID))
	h.hub.register <- conn

	// 启动读写协程
	go conn.writePump()
	go conn.readPump()
}

func (h *Handler) Hub() *Hub {
	return h.hub
}
